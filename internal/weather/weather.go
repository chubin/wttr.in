package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/renderer"
)

var ErrDataSource = errors.New("weather data source not available")

// Weatherer interface to fetch weather data based on location and language.
type Weatherer interface {
	GetWeather(lat, lon float64, lang string) ([]byte, error)
}

// IPLocator interface to fetch IP-related data.
type IPLocator interface {
	GetIPData(ip string) (*domain.IPData, error)
}

// Locator interface to fetch location-related data.
type Locator interface {
	GetLocation(location string) (*domain.Location, error)
}

// Renderer interface for rendering weather data into a visual representation.
type Renderer interface {
	Render(query domain.Query) (domain.RenderOutput, error)
}

// Formatter interface for converting rendered output into the final format.
type Formatter interface {
	Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error)
}

// QueryParser parses wttr.in / curl wttr.in style HTTP query strings
// and returns the result as a strongly-typed *options.Options struct.
type QueryParser interface {
	// Parse parses the raw query string (the part after the ? character)
	// and returns a populated *options.Options struct with all valid, active options set.
	//
	//   - Boolean flags without values are set to true (e.g. ?T -> Options.T = true)
	//   - Short flags can be bundled (e.g. ?0pq -> CurrentOnly=true, p=true, q=true)
	//   - Unknown, inactive, or invalid parameters cause an error
	//   - Validation rules from the YAML spec (ranges, regexps, allowed values, ...) are enforced
	//
	// If the query is empty (no ? or ? alone), a zero-valued *options.Options is typically returned
	// (all fields false/0/"").
	//
	// ctx can be used for cancellation, request-scoped logging, metrics collection, etc.
	// Most implementations will ignore it in the first version.
	Parse(context.Context, *http.Request, *options.Options) (*options.Options, error)

	// MustParse is a convenience variant that panics on error.
	// Mainly useful in tests, initialization code, or when invalid input is a programmer error.
	MustParse(context.Context, *http.Request, *options.Options) *options.Options
}

type RequestLogger interface {
	Log(r *http.Request) error
}

type UplinkProcessor interface {
	Route(opts *options.Options, r *http.Request, ipData *domain.IPData, location *domain.Location) (bool, *domain.CacheEntry, error)
}

// TimeTracker holds timing information for each step in the pipeline.
type TimeTracker struct {
	StepTimes []struct {
		Step string
		Time time.Duration
	}
}

func (tt *TimeTracker) Add(step string, t time.Duration) {
	tt.StepTimes = append(tt.StepTimes, struct {
		Step string
		Time time.Duration
	}{
		Step: step,
		Time: t,
	})
}

////////////////////////////////////////////////////////////////////////////////////////

// WeatherService struct holds the components necessary for processing a query.
type WeatherService struct {
	Weatherer       Weatherer
	Locator         Locator
	IPLocators      []IPLocator
	QueryParser     QueryParser
	Cacher          Cacher
	RequestLogger   RequestLogger
	UplinkProcessor UplinkProcessor
	RendererMap     map[string]Renderer
	FormatterMap    map[string]Formatter
}

// NewWeatherService initializes a new pipeline based on the provided options.
func NewWeatherService(
	weatherer Weatherer,
	locator Locator,
	ipLocators []IPLocator,
	queryParser QueryParser,
	cacher Cacher,
	requestLogger RequestLogger,
	uplinkProcessor UplinkProcessor,
	rendererMap map[string]Renderer,
	formatterMap map[string]Formatter,
) *WeatherService {
	return &WeatherService{
		Weatherer:       weatherer,
		Locator:         locator,
		IPLocators:      ipLocators,
		QueryParser:     queryParser,
		Cacher:          cacher,
		RequestLogger:   requestLogger,
		UplinkProcessor: uplinkProcessor,
		RendererMap:     rendererMap,
		FormatterMap:    formatterMap,
	}
}

// Helper to get the "best guess" client IP — handles proxies/load balancers
func getClientIP(r *http.Request) string {
	// Order of preference: X-Forwarded-For (first non-proxy IP), X-Real-IP, RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For can be "client-ip, proxy1, proxy2, ..."
		// Take the leftmost (original client) — but in production you should trust only known proxies
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// Fallback to direct connection (may be proxy/load-balancer IP)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr // worst case
}

// WeatherHandler handles weather requests with proper cache coalescing.
func (s *WeatherService) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	overallStart := time.Now()

	// Best-effort logging
	if err := s.RequestLogger.Log(r); err != nil {
		log.Println("RequestLogger error:", err)
	}

	bypassCache := r.URL.Query().Get("debug") != ""

	if bypassCache {
		s.serveFreshResponse(ctx, w, r)
		return
	}

	cacheKey := buildCacheKey(r)

	// Fast path: cache hit
	if entry := s.Cacher.Get(cacheKey); entry != nil {
		s.serveFromCache(w, entry)
		return
	}

	// Try to become the leader (atomic check + set)
	if !s.Cacher.SetInProgressIfNotExists(cacheKey) {
		// Someone else is already computing or we have a fresh entry now
		if entry, err := s.Cacher.WaitForCompletion(cacheKey, 12*time.Second); err == nil && entry != nil {
			s.serveFromCache(w, entry)
			return
		}
		// Timeout or failure → fall through and compute ourselves
	}

	// We are the leader → compute and store
	s.computeAndStore(ctx, w, r, cacheKey, overallStart)
}

// Helper to keep the main handler clean
func (s *WeatherService) serveFreshResponse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entry, err := s.computeResponse(ctx, r, &TimeTracker{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.serveFromCache(w, entry)
}

func (s *WeatherService) computeAndStore(ctx context.Context, w http.ResponseWriter, r *http.Request, cacheKey string, overallStart time.Time) {
	tracker := &TimeTracker{}

	defer func() {
		if rec := recover(); rec != nil {
			s.Cacher.Remove(cacheKey)
			panic(rec)
		}
	}()

	entry, err := s.computeResponse(ctx, r, tracker)
	if err != nil {
		s.Cacher.Remove(cacheKey)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success: store in cache
	s.Cacher.Set(cacheKey, *entry)

	tracker.Add("Total computation time", time.Since(overallStart))
	s.serveFromCache(w, entry)
}

func (s *WeatherService) serveFromCache(w http.ResponseWriter, e *domain.CacheEntry) {
	for k, vs := range e.Header {
		w.Header()[k] = vs
	}
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
	w.WriteHeader(e.StatusCode)
	w.Write(e.Body)
}

// computeResponse performs location resolution, weather fetch, rendering and formatting.
// Returns ready-to-cache domain.CacheEntry or error.
// Does NOT write to ResponseWriter — that stays in handler.
func (s *WeatherService) computeResponse(
	ctx context.Context, r *http.Request,
	tracker *TimeTracker,
) (*domain.CacheEntry, error) {
	debugRequested := r.URL.Query().Get("debug") != ""

	clientIP := getClientIP(r)

	// ── Determine location string ────────────────────────────────────────
	start := time.Now()
	path := r.URL.Path // cleanPath(r.URL.Path) // helper: trim / and extensions
	autoDetect := isAutoDetectPath(path)

	// 0. Locate IP
	var ipData *domain.IPData
	var errIP error
	for _, ipLocator := range s.IPLocators {
		ipData, errIP = ipLocator.GetIPData(clientIP)
		if errIP == nil && ipData != nil {
			break
		}
	}
	tracker.Add("IP locating", time.Since(start))

	// For testing purposes, it is possibe to simulate US-based clients
	// by setting country code in the HTTP headers of the request.
	if countryCode := r.Header.Get("X-Client-Country-Code"); countryCode != "" {
		ipData.CountryCode = countryCode
	}

	/////////////////////
	// This part should be moved to queryparser.
	ipOpts := options.Options{}
	if isClientInUSA(ipData) {
		ipOpts.UseMetric = false
		ipOpts.UseImperial = true
	} else {
		ipOpts.UseMetric = true
		ipOpts.UseImperial = false
	}
	///////////////////

	// 1. Parse options (cheap, always first)
	start = time.Now()
	opts, err := s.QueryParser.Parse(ctx, r, &ipOpts)
	if err != nil {
		return nil, err
	}
	tracker.Add("Options parsing", time.Since(start))

	locStr := opts.Location

	if ipData != nil {
		if autoDetect {
			if ipData.City == "" {
				locStr = fmt.Sprintf("%s,%s", ipData.Latitude, ipData.Longitude)
			} else {
				locStr = fmt.Sprintf("%s, %s, %s", ipData.City, ipData.Region, ipData.CountryCode)
				// 	log.Printf("Using old style location for: %s, %s, %s\n", ipData.City, ipData.Region, ipData.CountryCode)
				// 	locStr = ipData.City
			}
		}
	}
	tracker.Add("Determine location string + IP lookup", time.Since(start))

	if locStr == "" {
		// Temporary use Berlin as the default location
		locStr = "Berlin"
		log.Println("no location could be determined")
	}

	// ── Geocode ───────────────────────────────────────────────────────────
	start = time.Now()
	location, err := s.Locator.GetLocation(locStr)
	if err != nil && opts.View != "files" {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	tracker.Add("Geocode location", time.Since(start))

	// Here are we are ready to get the information from the upstream,
	// based on the information provided in the options.
	// We also have information about the IP and Geolocation,
	// which can be added to the headers.
	// ...
	var (
		formatOut *domain.FormatOutput
		query     domain.Query
	)

	// ── Build Query ───────────────────────────────────────────────────────
	start = time.Now()
	query = domain.Query{
		ClientData: &domain.ClientData{
			ClientIP:    clientIP,
			ClientAgent: r.UserAgent(),
		},
		Options:  opts,
		Location: location,
	}
	tracker.Add("Build Query object", time.Since(start))

	start = time.Now()
	var (
		isUpstream     bool
		uplinkResponse *domain.CacheEntry

		// compareWithUpstream is used when we the data
		// must be delivered both by the uplink, and the core
		// and then compared, e.g. for opts.View == 'line'.
		compareWithUpstream bool = false
	)
	isUpstream, uplinkResponse, err = s.UplinkProcessor.Route(opts, r, ipData, location)
	if isUpstream && !(compareWithUpstream && opts.View == "line") {
		if !debugRequested {
			return uplinkResponse, err
		}
		tracker.Add(fmt.Sprintf("Upstream processing (view=%s)", opts.View), time.Since(start))
	} else {
		// ── Fetch weather ─────────────────────────────────────────────────────
		start = time.Now()
		weatherBytes, err := s.Weatherer.GetWeather(location.Latitude, location.Longitude, opts.Lang)
		if err != nil {
			return nil, ErrDataSource
		}
		tracker.Add("Fetch weather data", time.Since(start))

		// ── Filling up Query ───────────────────────────────────────────────────────
		query.IPData = ipData
		query.Weather = (*domain.WeatherRaw)(&weatherBytes)

		// ── Render + Format ───────────────────────────────────────────────────
		start = time.Now()
		renderer := s.selectRenderer(opts.View)
		formatter := s.selectFormatter(opts.Output)

		renderOut, err := renderer.Render(query)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("render failed: %w", err)
		}

		formatOut, err = formatter.Format(&query, &renderOut)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("format failed: %w", err)
		}
		tracker.Add("Render + Format", time.Since(start))

		if isUpstream && (compareWithUpstream && opts.View == "line") {
			debugCompareOneLineRendering(opts.Format, string(uplinkResponse.Body), string(formatOut.Content))
			if !debugRequested {
				return uplinkResponse, err
			}
		}
	}

	if debugRequested {
		debugInfo := getDebugInfo(r, &query, locStr, tracker)
		return &domain.CacheEntry{
			Body:       []byte(debugInfo),
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type":  []string{"text/plain"},
				"Cache-Control": []string{"public, max-age=0"},
			},
		}, nil
	}

	return &domain.CacheEntry{
		Body:       formatOut.Content,
		StatusCode: http.StatusOK,
		Expires:    time.Now().Add(12 * time.Minute), // tune this TTL
		Header: http.Header{
			"Content-Type":  []string{formatOut.ContentType},
			"Cache-Control": []string{"public, max-age=600"},
		},
	}, nil
}

func isAutoDetectPath(p string) bool {
	return strings.Trim(p, "/") == ""
}

func getDebugInfo(
	r *http.Request,
	q *domain.Query,
	requestedLocStr string,
	timetracker *TimeTracker,
) string {
	var sb strings.Builder
	sb.WriteString("=== Weather Query Debug ===\n\n")

	sb.WriteString("Request Context:\n")
	sb.WriteString(fmt.Sprintf("  Method:       %s\n", r.Method))
	sb.WriteString(fmt.Sprintf("  Full URL:     %s\n", r.URL.String()))
	sb.WriteString(fmt.Sprintf("  Path:         %s\n", r.URL.Path))
	sb.WriteString(fmt.Sprintf("  Raw Query:    %s\n", r.URL.RawQuery))
	sb.WriteString(fmt.Sprintf("  Resolved IP:  %s\n", q.ClientData.ClientIP))
	sb.WriteString(fmt.Sprintf("  User-Agent:   %s\n", q.ClientData.ClientAgent))
	sb.WriteString("\n")

	sb.WriteString("Parsed Options:\n")
	sb.WriteString(prettyPrintOptions(q.Options))
	sb.WriteString("\n")

	sb.WriteString("Location Resolution:\n")
	sb.WriteString(fmt.Sprintf("  Requested loc string:  %q\n", requestedLocStr))
	if q.IPData != nil {
		sb.WriteString("  IP-derived data:\n")
		sb.WriteString(fmt.Sprintf("    IP:         %s\n", q.IPData.IP))
		sb.WriteString(fmt.Sprintf("    City:       %s\n", q.IPData.City))
		sb.WriteString(fmt.Sprintf("    Region:     %s\n", q.IPData.Region))
		sb.WriteString(fmt.Sprintf("    Country:    %s (%s)\n", q.IPData.Country, q.IPData.CountryCode))
		sb.WriteString(fmt.Sprintf("    Lat/Lon:    %s / %s\n", q.IPData.Latitude, q.IPData.Longitude))
	} else {
		sb.WriteString("  No IP lookup performed\n")
	}

	if q.Location != nil {
		sb.WriteString("\n  Final resolved location:\n")
		sb.WriteString(fmt.Sprintf("    Name:         %s\n", q.Location.Name))
		if q.Location.Country != "" {
			sb.WriteString(fmt.Sprintf("    Country:      %s (%s)\n", q.Location.Country, q.Location.CountryCode))
		}
		sb.WriteString(fmt.Sprintf("    Lat/Lon:      %.6f / %.6f\n", q.Location.Latitude, q.Location.Longitude))
		sb.WriteString(fmt.Sprintf("    Full addr:    %s\n", q.Location.FullAddress))
	} else {
		sb.WriteString("  No location resolved\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Weather Data:\n")
	if q.Weather != nil && len(*q.Weather) > 0 {
		sb.WriteString(fmt.Sprintf("  Fetched successfully (%d bytes)\n", len(*q.Weather)))
	} else {
		sb.WriteString("  Not fetched (no valid location/used uplink)\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Time Tracking:\n")
	for _, timeStep := range timetracker.StepTimes {
		sb.WriteString(fmt.Sprintf("  %s: %v\n", timeStep.Step, timeStep.Time))
	}

	return sb.String()
}

func prettyPrintOptions(o *options.Options) string {
	if o == nil {
		return "  (nil)\n"
	}

	data, err := json.MarshalIndent(o, "  ", "  ")
	if err != nil {
		return fmt.Sprintf("  Error: %v\n", err)
	}

	return "  " + string(data) + "\n"
}

// selectRenderer chooses the appropriate renderer based on the format option.
func (s *WeatherService) selectRenderer(view string) Renderer {
	if rndrer, found := s.RendererMap[view]; found {
		return rndrer
	} else {
		log.Println("Unknown renderer for view: ", view)
		return &renderer.V1Renderer{} // If no format specified, use v1 renderer
	}
}

// selectFormatter chooses the appropriate formatter based on the output format option.
func (s *WeatherService) selectFormatter(output string) Formatter {
	switch output {
	case "terminal", "text", "ansi":
		return s.FormatterMap["text"]
	case "html":
		return s.FormatterMap["html"]
	case "png":
		return s.FormatterMap["png"]
	default:
		return s.FormatterMap["text"]
	}
}

func debugCompareOneLineRendering(format string, uplinkResponse string, internalResponse string) {
	if uplinkResponse != internalResponse {
		// Prepare the log message using the provided format and responses
		logMessage := fmt.Sprintf("---\nFORMAT:\n%s\n###########\nUPLINK:\n%s\nINTERNAL:\n%s\n^^^\n", format, uplinkResponse, internalResponse)

		// Open the log file in append mode, create if it doesn't exist
		file, err := os.OpenFile("/tmp/oneline-rendering-errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// If there's an error opening the file, log it to stderr
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
			return
		}
		defer file.Close()

		// Write the log message to the file
		if _, err := file.WriteString(logMessage + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
		}
	}
}

func isClientInUSA(ipData *domain.IPData) bool {
	return ipData.CountryCode == "US"
}
