package weather

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// WWOConfig holds configuration for World Weather Online (or similar) API
type WWOConfig struct {
	BaseURL  string `yaml:"baseUrl"`
	Key      string `yaml:"key"`
	MaxConns int    `yaml:"maxConns,omitempty"`
}

// WeatherClient fetches weather data with controlled concurrency using a semaphore.
type WeatherClient struct {
	baseURL  string
	sem      chan struct{} // semaphore to limit concurrent connections
	maxConns int
}

// ErrNoFreeConnection is returned when all connection slots are occupied.
var ErrNoFreeConnection = fmt.Errorf("heating up, try again later")

// NewWeatherClient creates a new WeatherClient.
func NewWeatherClient(cfg *WWOConfig) *WeatherClient {
	if cfg == nil {
		panic("weather config is nil")
	}
	if cfg.BaseURL == "" {
		panic("empty baseURL in weather config")
	}
	if cfg.Key == "" {
		panic("missing/empty API key in weather config")
	}

	maxConns := cfg.MaxConns
	if maxConns < 1 {
		maxConns = 100 // sensible default
	}

	// Replace {key} once at client creation
	baseURL := strings.ReplaceAll(cfg.BaseURL, "{key}", cfg.Key)

	return &WeatherClient{
		baseURL:  baseURL,
		sem:      make(chan struct{}, maxConns),
		maxConns: maxConns,
	}
}

// GetWeather fetches weather data for given coordinates while respecting the concurrency limit.
// Returns ErrNoFreeConnection immediately if no slot is available (non-blocking).
func (wc *WeatherClient) GetWeather(lat, lon float64, lang string) ([]byte, error) {
	// Non-blocking acquire
	select {
	case wc.sem <- struct{}{}:
		// Slot acquired - ensure we release it when done
		defer func() { <-wc.sem }()
	default:
		return nil, ErrNoFreeConnection
	}

	// Build URL with placeholders
	url := wc.buildURL(lat, lon, lang)

	logrus.WithFields(logrus.Fields{
		"lat":  lat,
		"lon":  lon,
		"lang": lang,
	}).Debug("Accessing weather API")

	// HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read weather response body: %w", err)
	}

	return body, nil
}

// buildURL constructs the final URL by replacing placeholders.
// This is extracted for better testability and clarity.
func (wc *WeatherClient) buildURL(lat, lon float64, lang string) string {
	url := strings.ReplaceAll(wc.baseURL, "{lat}", fmt.Sprintf("%.6f", lat))
	url = strings.ReplaceAll(url, "{lon}", fmt.Sprintf("%.6f", lon))
	url = strings.ReplaceAll(url, "{lang}", lang)
	return url
}