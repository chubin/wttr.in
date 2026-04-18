package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// WWOConfig holds configuration for the weather API (e.g. World Weather Online)
type WWOConfig struct {
	BaseURL  string `yaml:"baseUrl"`
	Key      string `yaml:"key"`
	MaxConns int    `yaml:"maxConns,omitempty"`
}

// WeatherClient fetches weather data with controlled concurrency and connection reuse.
type WeatherClient struct {
	baseURL    string
	sem        chan struct{} // semaphore to limit concurrent requests
	maxConns   int
	httpClient *http.Client // reused HTTP client with connection pooling
}

// ErrNoFreeConnection is returned when the maximum number of parallel connections is reached.
var ErrNoFreeConnection = fmt.Errorf("heating up, try again later")

// NewWeatherClient creates a new WeatherClient with proper connection pooling.
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
		maxConns = 100 // reasonable default
	}

	// Replace API key once during initialization
	baseURL := strings.ReplaceAll(cfg.BaseURL, "{key}", cfg.Key)

	// Configure HTTP transport for connection reuse and pooling
	transport := &http.Transport{
		MaxIdleConns:        100,              // total idle connections across all hosts
		MaxIdleConnsPerHost: 30,               // idle connections to the weather API host
		MaxConnsPerHost:     maxConns,         // maximum concurrent connections to the host
		IdleConnTimeout:     90 * time.Second, // how long to keep idle connections alive
	}

	return &WeatherClient{
		baseURL:  baseURL,
		sem:      make(chan struct{}, maxConns),
		maxConns: maxConns,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

// GetWeather fetches weather data while respecting the concurrency limit.
// It reuses HTTP connections thanks to the shared http.Client.
func (wc *WeatherClient) GetWeather(lat, lon float64, lang string) ([]byte, error) {
	// Non-blocking attempt to acquire a semaphore slot
	select {
	case wc.sem <- struct{}{}:
		// Slot acquired - release when function returns
		defer func() { <-wc.sem }()
	default:
		// All slots are busy
		return nil, ErrNoFreeConnection
	}

	// Build the final URL
	url := wc.buildURL(lat, lon, lang)

	logrus.WithFields(logrus.Fields{
		"lat":  lat,
		"lon":  lon,
		"lang": lang,
	}).Debug("Accessing weather API")

	// Use the shared http.Client (this enables connection reuse)
	resp, err := wc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read weather response body: %w", err)
	}

	////////////////
	// Remove 'data' wrapper.
	var data struct {
		Data interface{} `json:"data"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid data format")
	}

	dataBytes, err := json.MarshalIndent(data.Data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("invalid data format")
	}
	//////////////////

	return dataBytes, nil
}

// buildURL replaces placeholders in the base URL with actual values.
func (wc *WeatherClient) buildURL(lat, lon float64, lang string) string {
	url := strings.ReplaceAll(wc.baseURL, "{lat}", fmt.Sprintf("%.6f", lat))
	url = strings.ReplaceAll(url, "{lon}", fmt.Sprintf("%.6f", lon))
	url = strings.ReplaceAll(url, "{lang}", lang)
	return url
}

// Close closes idle HTTP connections. Call this when shutting down the application.
func (wc *WeatherClient) Close() {
	if wc.httpClient != nil {
		if transport, ok := wc.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}
