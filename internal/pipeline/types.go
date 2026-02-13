package pipeline

import (
	"net/http"

	"github.com/chubin/wttr.go/internal/query"
)

// Weatherer interface to fetch weather data based on location and language.
type Weatherer interface {
	GetWeather(lat, lon float64, lang string) ([]byte, error)
}

// IPDataer interface to fetch IP-related data.
type IPDataer interface {
	GetIPData(ip string) (IPData, error)
}

// Renderer interface for rendering weather data into a visual representation.
type Renderer interface {
	Render(query Query) (RenderOutput, error)
}

// Formatter interface for converting rendered output into the final format.
type Formatter interface {
	Format(output RenderOutput) (FormatOutput, error)
}

// ClientData holds information about the client making the request.
type ClientData struct {
	ClientIP    string
	ClientAgent string
}

// IPData holds information about the client's IP address and location.
type IPData struct {
	IP           string
	LocationAddr string
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

// Pipeline struct holds the components necessary for processing a query.
type Pipeline struct {
	Weatherer Weatherer
	IPDataer  IPDataer
	Renderer  Renderer
	Formatter Formatter
}

// NewPipeline initializes a new pipeline based on the provided options.
func NewPipeline(opts *query.Options, weatherer Weatherer, ipDataer IPDataer) *Pipeline {
	renderer := selectRenderer(opts.Format)
	formatter := selectFormatter(opts.Format)
	return &Pipeline{
		Weatherer: weatherer,
		IPDataer:  ipDataer,
		Renderer:  renderer,
		Formatter: formatter,
	}
}

// Process runs the pipeline to process the query and produce the final output.
func (p *Pipeline) Process(query Query) (FormatOutput, error) {
	// Step 1: Render the query data into an intermediate format (ANSI)
	renderOutput, err := p.Renderer.Render(query)
	if err != nil {
		return FormatOutput{}, err
	}

	// Step 2: Format the rendered output into the final format
	formatOutput, err := p.Formatter.Format(renderOutput)
	if err != nil {
		return FormatOutput{}, err
	}

	return formatOutput, nil
}

// HTTP Handler for processing incoming weather queries.
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Parse incoming HTTP query into Options
	opts, err := parseQueryOptions(r)
	if err != nil {
		http.Error(w, "Failed to parse query options", http.StatusBadRequest)
		return
	}

	// Initialize pipeline with the parsed options
	pipeline := NewPipeline(opts, nil, nil) // Weatherer and IPDataer to be injected or initialized

	// Build the query struct with client data and options
	query := Query{
		ClientData: &ClientData{
			ClientIP:    r.RemoteAddr,
			ClientAgent: r.UserAgent(),
		},
		Options: opts,
	}

	// Process the query through the pipeline
	output, err := pipeline.Process(query)
	if err != nil {
		http.Error(w, "Failed to process query", http.StatusInternalServerError)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", output.ContentType)

	// Write the formatted output to the HTTP response
	_, err = w.Write(output.Content)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
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

// Weatherer and IPDataer implementations are also stubs to be provided externally.
// They should be injected into NewPipeline during initialization.
