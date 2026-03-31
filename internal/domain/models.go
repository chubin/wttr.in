// Package domain contains the core domain models and data structures
// of the weather service.
//
// This package represents the central, stable layer of the application.
// It defines the main business entities and DTOs that are shared across
// multiple packages (weather, renderer, location, ip, query, etc.).
//
// All types in this package are designed to be:
//   - Independent of any external framework, database, or transport layer
//   - Immutable where appropriate
//   - Used as the canonical data representation throughout the service
//
// Following clean architecture principles, higher-level packages depend
// on domain, but domain depends on nothing else (except for a few small
// standard library types and the query.Options struct).
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

// IPData holds geolocation information derived from the client's IP address.
type IPData struct {
	IP          string
	CountryCode string
	Country     string
	Region      string
	City        string
	Latitude    string
	Longitude   string
}

// Location represents a geographic location for which weather data is requested.
//
// It can be populated either from user input (city name, coordinates) or
// from IP-based geolocation.
type Location struct {
	Name        string
	Country     string
	CountryCode string
	Latitude    float64
	Longitude   float64
	FullAddress string
	TimeZone    string
}

// CacheEntry represents a cached response entry.
//
// It is designed to be immutable once created and stored in any cache
// implementation. The entry includes everything needed to reconstruct
// a proper HTTP response.
type CacheEntry struct {
	// Body is the raw response body bytes
	Body []byte

	// Header contains the HTTP response headers to return to the client
	Header http.Header

	// StatusCode is the HTTP status code (e.g. 200, 404, 429)
	StatusCode int

	// Expires is the absolute time after which this cache entry
	// should be considered stale and no longer served.
	Expires time.Time
}

// WeatherRaw is a type alias for raw weather data returned by a weather backend.
//
// It is typically JSON or another serialized format. The actual structure
// is handled by specific renderers and parsers.
type WeatherRaw []byte

// RenderOutput is the intermediate result produced by a Renderer.
//
// It typically contains ANSI-colored or plain text output intended for
// terminal or further processing.
type RenderOutput struct {
	Content []byte
}

// FormatOutput represents the final output ready to be sent to the client.
//
// It includes both the content and the appropriate HTTP Content-Type header.
type FormatOutput struct {
	Content     []byte
	ContentType string
}

// Query represents the complete context of a single weather request.
//
// It aggregates all information needed to process and respond to a request:
// client details, parsed options, geolocation data, resolved location,
// and the fetched weather data (once available).
type Query struct {
	ClientData *ClientData
	Options    *query.Options
	IPData     *IPData
	Location   *Location
	Weather    *WeatherRaw
}
