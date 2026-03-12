package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/chubin/wttr.go/internal/query"
)

// Weatherer interface to fetch weather data based on location and language.
type Weatherer interface {
	GetWeather(lat, lon float64, lang string) ([]byte, error)
}

// IPLocator interface to fetch IP-related data.
type IPLocator interface {
	GetIPData(ip string) (*IPData, error)
}

// Locator interface to fetch location-related data.
type Locator interface {
	GetLocation(location string) (*Location, error)
}

// Renderer interface for rendering weather data into a visual representation.
type Renderer interface {
	Render(query Query) (RenderOutput, error)
}

// Formatter interface for converting rendered output into the final format.
type Formatter interface {
	Format(output RenderOutput) (*FormatOutput, error)
}

// QueryParser parses wttr.in / curl wttr.in style HTTP query strings
// and returns the result as a strongly-typed *query.Options struct.
type QueryParser interface {
	// Parse parses the raw query string (the part after the ? character)
	// and returns a populated *query.Options struct with all valid, active options set.
	//
	//   - Boolean flags without values are set to true (e.g. ?T -> Options.T = true)
	//   - Short flags can be bundled (e.g. ?0pq -> CurrentOnly=true, p=true, q=true)
	//   - Unknown, inactive, or invalid parameters cause an error
	//   - Validation rules from the YAML spec (ranges, regexps, allowed values, ...) are enforced
	//
	// If the query is empty (no ? or ? alone), a zero-valued *query.Options is typically returned
	// (all fields false/0/"").
	//
	// ctx can be used for cancellation, request-scoped logging, metrics collection, etc.
	// Most implementations will ignore it in the first version.
	Parse(ctx context.Context, r *http.Request) (*query.Options, error)

	// MustParse is a convenience variant that panics on error.
	// Mainly useful in tests, initialization code, or when invalid input is a programmer error.
	MustParse(ctx context.Context, r *http.Request) *query.Options
}

type RequestLogger interface {
	Log(r *http.Request) error
}

type UplinkProcessor interface {
	Route(opts *query.Options, r *http.Request, ipData *IPData, location *Location) (bool, *CacheEntry, error)
}

// ClientData holds information about the client making the request.
type ClientData struct {
	ClientIP    string
	ClientAgent string
}

// IPData holds information about the client's IP address and location.
type IPData struct {
	IP          string
	CountryCode string
	Country     string
	Region      string
	City        string
	Latitude    string
	Longitude   string
}

// Location holds detailed information about a specific location.
type Location struct {
	Name        string
	Country     string
	CountryCode string
	Latitude    float64
	Longitude   float64
	FullAddress string
}

// WeatherData represents the internal weather data as raw bytes.
type WeatherData []byte

// Query holds all data necessary for processing a weather query.
type Query struct {
	ClientData *ClientData
	Options    *query.Options
	IPData     *IPData
	Location   *Location
	Weather    *WeatherData
}

// RenderOutput represents the intermediate output from a renderer (ANSI format).
type RenderOutput struct {
	Content []byte
}

// FormatOutput represents the final formatted output to be sent to the client.
type FormatOutput struct {
	Content     []byte
	ContentType string
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
) *WeatherService {
	return &WeatherService{
		Weatherer:       weatherer,
		Locator:         locator,
		IPLocators:      ipLocators,
		QueryParser:     queryParser,
		Cacher:          cacher,
		RequestLogger:   requestLogger,
		UplinkProcessor: uplinkProcessor,
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

// WeatherHandler is now much shorter — mainly orchestration + caching
func (s *WeatherService) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var (
		bypassCache bool
		cacheKey    string
		entry       *CacheEntry
		err         error
	)

	overallStart := time.Now()
	tracker := TimeTracker{}

	// Log incoming request.
	// If the logging was not successful, write a warning and continue.
	if err := s.RequestLogger.Log(r); err != nil {
		log.Println(err)
	}

	if r.URL.Query().Get("debug") != "" {
		bypassCache = true
	}

	if !bypassCache {
		// 3. Build cache key (must include everything that affects output)
		cacheKey = buildCacheKey(r)

		// 4. Fast path: cache hit
		if entry := s.Cacher.Get(cacheKey); entry != nil {
			s.serveFromCache(w, entry)
			return
		}

		// 5. Coalescing: if someone else is already computing → wait
		if s.Cacher.IsInProgress(cacheKey) {
			entry, err := s.Cacher.WaitForCompletion(cacheKey, 12*time.Second)
			if err == nil && entry != nil {
				s.serveFromCache(w, entry)
				return
			}
			// timeout → we'll compute it ourselves
		}

		// 6. Become the leader: mark in-progress and compute
		s.Cacher.SetInProgress(cacheKey)

		// Important: recover from panic to clean up in-progress flag
		defer func() {
			if rec := recover(); rec != nil {
				s.Cacher.Remove(cacheKey)
				panic(rec)
			}
		}()
	}

	// 7. The heavy part — now extracted
	entry, err = s.computeResponse(ctx, r, &tracker)
	if err != nil {
		s.Cacher.Remove(cacheKey)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// This information will not be visible,
	// because the debug output is already in response.
	tracker.Add("Total computation time", time.Since(overallStart))

	if !bypassCache {
		// 8. Store successful result
		s.Cacher.Set(cacheKey, *entry)
	}

	// 9. Send to client
	s.serveFromCache(w, entry) // reuse same helper
}

func (s *WeatherService) serveFromCache(w http.ResponseWriter, e *CacheEntry) {
	for k, vs := range e.Header {
		w.Header()[k] = vs
	}
	w.WriteHeader(e.StatusCode)
	w.Write(e.Body)
}

// computeResponse performs location resolution, weather fetch, rendering and formatting.
// Returns ready-to-cache CacheEntry or error.
// Does NOT write to ResponseWriter — that stays in handler.
func (s *WeatherService) computeResponse(
	ctx context.Context, r *http.Request,
	tracker *TimeTracker,
) (*CacheEntry, error) {
	debugRequested := r.URL.Query().Get("debug") != ""

	clientIP := getClientIP(r)

	// ── Determine location string ────────────────────────────────────────
	start := time.Now()
	path := r.URL.Path // cleanPath(r.URL.Path) // helper: trim / and extensions
	autoDetect := isAutoDetectPath(path)

	// 1. Parse options (cheap, always first)
	opts, err := s.QueryParser.Parse(ctx, r)
	if err != nil {
		return nil, err
	}
	tracker.Add("Options parsing", time.Since(start))

	locStr := opts.Location
	var ipData *IPData
	if autoDetect {
		var errIP error
		for _, ipLocator := range s.IPLocators {
			ipData, errIP = ipLocator.GetIPData(clientIP)
			if errIP == nil && ipData != nil {
				locStr = ipData.City
				if locStr == "" {
					locStr = fmt.Sprintf("%s,%s", ipData.Latitude, ipData.Longitude)
				}
				break
			}
		}
		tracker.Add("Determine location string + IP lookup", time.Since(start))
	}

	if locStr == "" {
		return nil, fmt.Errorf("no location could be determined")
	}

	// ── Geocode ───────────────────────────────────────────────────────────
	start = time.Now()
	location, err := s.Locator.GetLocation(locStr)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	tracker.Add("Geocode location", time.Since(start))

	// Here are we are ready to get the information from the upstream,
	// based on the information provided in the options.
	// We also have information about the IP and Geolocation,
	// which can be added to the headers.
	// ...
	isUpstream, uplinkResponse, err := s.UplinkProcessor.Route(opts, r, ipData, location)
	if isUpstream {
		return uplinkResponse, err
	}

	// ── Fetch weather ─────────────────────────────────────────────────────
	start = time.Now()
	weatherBytes, err := s.Weatherer.GetWeather(location.Latitude, location.Longitude, opts.Lang)
	if err != nil {
		return nil, fmt.Errorf("weather fetch failed: %w", err)
	}
	tracker.Add("Fetch weather data", time.Since(start))

	// ── Build Query ───────────────────────────────────────────────────────
	start = time.Now()
	query := Query{
		ClientData: &ClientData{
			ClientIP:    clientIP,
			ClientAgent: r.UserAgent(),
		},
		Options:  opts,
		IPData:   ipData,
		Location: location,
		Weather:  (*WeatherData)(&weatherBytes),
	}
	tracker.Add("Build Query object", time.Since(start))

	// ── Render + Format ───────────────────────────────────────────────────
	start = time.Now()
	renderer := selectRenderer(opts.Format)
	formatter := selectFormatter(opts.Format)

	renderOut, err := renderer.Render(query)
	if err != nil {
		return nil, fmt.Errorf("render failed: %w", err)
	}

	formatOut, err := formatter.Format(renderOut)
	if err != nil {
		return nil, fmt.Errorf("format failed: %w", err)
	}
	tracker.Add("Render + Format", time.Since(start))

	if debugRequested {
		debugInfo := getDebugInfo(r, &query, locStr, tracker)
		formatOut = &FormatOutput{
			Content:     []byte(debugInfo),
			ContentType: "text/plain",
		}
	}

	return &CacheEntry{
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
	q *Query,
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
	if len(*q.Weather) > 0 {
		sb.WriteString(fmt.Sprintf("  Fetched successfully (%d bytes)\n", len(*q.Weather)))
	} else {
		sb.WriteString("  Not fetched (no valid location)\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Time Tracking:\n")
	for _, timeStep := range timetracker.StepTimes {
		sb.WriteString(fmt.Sprintf("  %s: %v\n", timeStep.Step, timeStep.Time))
	}

	return sb.String()
}

func prettyPrintOptions(o *query.Options) string {
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
func selectRenderer(format string) Renderer {
	switch format {
	case "":
		return &V1Renderer{} // If no format specified, use v1 renderer
	case "v1":
		return &V1Renderer{}
	case "v2":
		return &V2Renderer{}
	case "j1":
		return &J1Renderer{}
	case "j2":
		return &J2Renderer{}
	default:
		return &OnelineRenderer{} // Default to oneline renderer
	}
}

// selectFormatter chooses the appropriate formatter based on the format option.
func selectFormatter(format string) Formatter {
	if format == "j1" || format == "j2" {
		format = "json"
	}

	switch format {
	case "terminal":
		return &TerminalFormatter{}
	case "browser":
		return &BrowserFormatter{}
	case "png":
		return &PNGFormatter{}
	case "json":
		return &JSONFormatter{}
	default:
		return &TerminalFormatter{} // Default to terminal formatter
	}
}

// Renderer Implementations (Stubs)
type V1Renderer struct{}

func (r *V1Renderer) Render(query Query) (RenderOutput, error) {
	// Stub: To be implemented
	return RenderOutput{}, nil
}

type V2Renderer struct{}

func (r *V2Renderer) Render(query Query) (RenderOutput, error) {
	// Stub: To be implemented
	return RenderOutput{}, nil
}

type OnelineRenderer struct{}

func (r *OnelineRenderer) Render(query Query) (RenderOutput, error) {
	// Stub: To be implemented
	return RenderOutput{}, nil
}

// Formatter Implementations (Stubs)
type TerminalFormatter struct{}

func (f *TerminalFormatter) Format(output RenderOutput) (*FormatOutput, error) {
	// Stub: To be implemented
	return &FormatOutput{}, nil
}

type BrowserFormatter struct{}

func (f *BrowserFormatter) Format(output RenderOutput) (*FormatOutput, error) {
	// Stub: To be implemented
	return &FormatOutput{}, nil
}

type PNGFormatter struct{}

func (f *PNGFormatter) Format(output RenderOutput) (*FormatOutput, error) {
	// Stub: To be implemented
	return &FormatOutput{}, nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(output RenderOutput) (*FormatOutput, error) {
	return &FormatOutput{
		Content:     output.Content,
		ContentType: "application/json",
	}, nil
}

// Stub Functions to be Implemented Separately
// parseQueryOptions parses the incoming HTTP request into Options.
func parseQueryOptions(r *http.Request) (*query.Options, error) {
	// Stub: To be implemented
	return nil, nil
}

// Weatherer and IPLocator implementations are also stubs to be provided externally.
// They should be injected into NewWeatherService during initialization.
