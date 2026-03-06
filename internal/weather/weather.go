package weather

import (
	"context"
	"encoding/json"
	"fmt"
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
	Format(output RenderOutput) (FormatOutput, error)
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
	Parse(ctx context.Context, query string) (*query.Options, error)

	// MustParse is a convenience variant that panics on error.
	// Mainly useful in tests, initialization code, or when invalid input is a programmer error.
	MustParse(ctx context.Context, query string) *query.Options
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

////////////////////////////////////////////////////////////////////////////////////////

// WeatherService struct holds the components necessary for processing a query.
type WeatherService struct {
	Weatherer   Weatherer
	Locator     Locator
	IPLocator   IPLocator
	QueryParser QueryParser
}

// NewWeatherService initializes a new pipeline based on the provided options.
func NewWeatherService(weatherer Weatherer, locator Locator, ipLocator IPLocator, queryParser QueryParser) *WeatherService {
	return &WeatherService{
		Weatherer:   weatherer,
		Locator:     locator,
		IPLocator:   ipLocator,
		QueryParser: queryParser,
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

// WeatherHandler processes incoming weather queries.
func (s *WeatherService) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	timetracker := TimeTracker{}

	ctx := r.Context()
	start := time.Now()
	opts, err := s.QueryParser.Parse(ctx, r.URL.RawQuery)
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Parse query options", time.Since(start)})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse query options: %s", err), http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)

	// ── 1. Determine location string from path ───────────────────────────────
	start = time.Now()
	path := strings.Trim(r.URL.Path, "/")
	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") {
		path = strings.TrimSuffix(path, ".png")
		path = strings.TrimSuffix(path, ".jpg")
	}
	// Special cases like /:help, moon, etc. can be added later

	var locStr string
	autoDetect := (path == "" || path == ":help" || path == "help") // add more special paths if needed
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Determine location string from path", time.Since(start)})

	// ── 2. Get IP-based data (always — useful for fallback + logging) ────────
	start = time.Now()
	var ipData *IPData
	var errIP error
	if autoDetect || path == "" {
		ipData, errIP = s.IPLocator.GetIPData(clientIP)
		if errIP == nil && ipData != nil {
			locStr = ipData.City
			if locStr == "" {
				// fallback to coords if city missing
				locStr = fmt.Sprintf("%s,%s", ipData.Latitude, ipData.Longitude)
			}
		}
		// If IP lookup fails → locStr stays "", will error later
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Get IP-based data", time.Since(start)})

	// ── 3. Override with explicit path location if provided ──────────────────
	start = time.Now()
	if !autoDetect && path != "" {
		locStr = path
		// Optional: reset ipData to nil or keep for logging
		ipData = nil
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Override with explicit path location if provided", time.Since(start)})

	// ── 4. Resolve final Location (geocode locStr) ───────────────────────────
	start = time.Now()
	var location *Location
	var errLoc error
	if locStr != "" {
		location, errLoc = s.Locator.GetLocation(locStr)
		// if errLoc != nil {
		// 	http.Error(w, fmt.Sprintf("Location not found: %s", locStr), http.StatusNotFound)
		// 	return
		// }
	} else {
		// No location at all (very rare — IP failed and no path)
		http.Error(w, "No location provided and IP detection failed", http.StatusBadRequest)
		return
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Resolve final Location (geocode locStr)", time.Since(start)})

	// ── 5. Fetch weather data ────────────────────────────────────────────────
	start = time.Now()
	var weatherBytes []byte
	var errWeather error
	if location != nil {
		weatherBytes, errWeather = s.Weatherer.GetWeather(location.Latitude, location.Longitude, opts.Lang)
		if errWeather != nil {
			// // In production: consider fallback / cached / different provider
			// http.Error(w, "Failed to retrieve weather data", http.StatusBadGateway)
			// return
		}
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Fetch weather data", time.Since(start)})

	// ── 6. Build complete Query ──────────────────────────────────────────────
	start = time.Now()
	query := Query{
		ClientData: &ClientData{
			ClientIP:    clientIP,
			ClientAgent: r.UserAgent(),
		},
		Options:  opts,
		IPData:   ipData, // may be nil if explicit location was used
		Location: location,
		Weather:  (*WeatherData)(&weatherBytes),
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Build complete Query", time.Since(start)})

	// ── Debug check – AFTER everything is resolved ───────────────────────────
	debugRequested := opts.Debug ||
		r.URL.Query().Get("debug") != "" ||
		r.Header.Get("X-Debug") == "1"

	if debugRequested {
		s.serveDebugInfo(w, r, &query, locStr, errIP, errLoc, errWeather, &timetracker)
		return
	}

	// ── Renderer + Formatter pipeline (unchanged from previous step) ─────────
	start = time.Now()
	renderer := selectRenderer(opts.Format)
	formatter := selectFormatter(opts.Format)

	renderOutput, err := renderer.Render(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rendering failed: %v", err), http.StatusInternalServerError)
		return
	}

	formatOutput, err := formatter.Format(renderOutput)
	if err != nil {
		http.Error(w, fmt.Sprintf("Formatting failed: %v", err), http.StatusInternalServerError)
		return
	}
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Renderer and Formatter pipeline", time.Since(start)})

	w.Header().Set("Content-Type", formatOutput.ContentType)
	// Optional: short cache — weather usually stable for 5–15 min
	w.Header().Set("Cache-Control", "public, max-age=600")

	_, err = w.Write(formatOutput.Content)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	totalDuration := time.Since(startTime)
	timetracker.StepTimes = append(timetracker.StepTimes, struct {
		Step string
		Time time.Duration
	}{"Total", totalDuration})
}

func (s *WeatherService) serveDebugInfo(
	w http.ResponseWriter,
	r *http.Request,
	q *Query,
	requestedLocStr string,
	ipErr, locErr, weatherErr error,
	timetracker *TimeTracker,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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
	} else if ipErr != nil {
		sb.WriteString(fmt.Sprintf("  IP lookup error: %v\n", ipErr))
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
	} else if locErr != nil {
		sb.WriteString(fmt.Sprintf("  Geocoding error: %v\n", locErr))
	} else {
		sb.WriteString("  No location resolved\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Weather Data:\n")
	if len(*q.Weather) > 0 {
		sb.WriteString(fmt.Sprintf("  Fetched successfully (%d bytes)\n", len(*q.Weather)))
	} else if weatherErr != nil {
		sb.WriteString(fmt.Sprintf("  Fetch error: %v\n", weatherErr))
	} else {
		sb.WriteString("  Not fetched (no valid location)\n")
	}
	sb.WriteString("\n")

	sb.WriteString("Time Tracking:\n")
	for _, timeStep := range timetracker.StepTimes {
		sb.WriteString(fmt.Sprintf("  %s: %v\n", timeStep.Step, timeStep.Time))
	}

	fmt.Fprint(w, sb.String())
}

func prettyPrintOptions(o *query.Options) string {
	if o == nil {
		return "  (nil)\n"
	}

	data, err := json.MarshalIndent(o, "  ", "  ")
	if err != nil {
		return fmt.Sprintf("  Error: %v\n", err)
	}

	return string(data) + "\n"
}

// selectRenderer chooses the appropriate renderer based on the format option.
func selectRenderer(format string) Renderer {
	switch format {
	case "v1":
		return &V1Renderer{}
	case "v2":
		return &V2Renderer{}
	case "oneline":
		return &OnelineRenderer{}
	case "j1":
		return &J1Renderer{}
	case "j2":
		return &J2Renderer{}
	default:
		return &V1Renderer{} // Default to v1 renderer
	}
}

// selectFormatter chooses the appropriate formatter based on the format option.
func selectFormatter(format string) Formatter {
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

func (f *TerminalFormatter) Format(output RenderOutput) (FormatOutput, error) {
	// Stub: To be implemented
	return FormatOutput{}, nil
}

type BrowserFormatter struct{}

func (f *BrowserFormatter) Format(output RenderOutput) (FormatOutput, error) {
	// Stub: To be implemented
	return FormatOutput{}, nil
}

type PNGFormatter struct{}

func (f *PNGFormatter) Format(output RenderOutput) (FormatOutput, error) {
	// Stub: To be implemented
	return FormatOutput{}, nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(output RenderOutput) (FormatOutput, error) {
	// Stub: To be implemented
	return FormatOutput{}, nil
}

// Stub Functions to be Implemented Separately
// parseQueryOptions parses the incoming HTTP request into Options.
func parseQueryOptions(r *http.Request) (*query.Options, error) {
	// Stub: To be implemented
	return nil, nil
}

// Weatherer and IPLocator implementations are also stubs to be provided externally.
// They should be injected into NewWeatherService during initialization.
