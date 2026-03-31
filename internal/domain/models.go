package domain

import (
	"net/http"
	"time"

	"github.com/chubin/wttr.in/internal/query"
)

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
	TimeZone    string
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

// CacheEntry represents a cached HTTP response.
// It is immutable once stored in the cache.
type CacheEntry struct {
	// Body is the response body bytes
	Body []byte

	// Header contains the HTTP response headers to return to the client
	Header http.Header

	// StatusCode is the HTTP status code (200, 404, etc.)
	StatusCode int

	// Expires is the absolute time after which this entry should be considered stale
	Expires time.Time
}
