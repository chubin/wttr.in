package weather

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type WWOConfig struct {
	BaseURL  string `yaml:"baseUrl"`
	Key      string `yaml:"key"`
	MaxConns int    `yaml:"maxConns,omitempty"`
}

// WeatherClient fetches weather data with controlled concurrency
type WeatherClient struct {
	baseURL  string
	sem      chan struct{} // semaphore to limit concurrent connections
	maxConns int
}

// NewWeatherClient creates a new WeatherClient with a maximum number of parallel connections
func NewWeatherClient(cfg *WWOConfig) *WeatherClient {
	if cfg.BaseURL == "" {
		panic("empty baseURL")
	}
	if cfg.Key == "" {
		panic("missing/empty key")
	}

	maxConns := cfg.MaxConns
	if maxConns < 1 {
		maxConns = 100 // reasonable default
	}

	baseURL := cfg.BaseURL
	baseURL = strings.Replace(baseURL, "{key}", cfg.Key, 1)

	return &WeatherClient{
		baseURL:  baseURL,
		sem:      make(chan struct{}, maxConns),
		maxConns: maxConns,
	}
}

// GetWeather fetches weather data while respecting the maximum number of parallel connections.
// Returns ErrNoFreeConnection if all slots are currently occupied.
func (wc *WeatherClient) GetWeather(lat, lon float64, lang string) ([]byte, error) {
	// Try to acquire a connection slot (non-blocking)
	select {
	case wc.sem <- struct{}{}:
		// Slot acquired successfully
		defer func() { <-wc.sem }() // release slot when done
	default:
		// No free connection slot available
		return nil, ErrNoFreeConnection
	}

	// Replace placeholders
	url := strings.Replace(wc.baseURL, "{lat}", fmt.Sprintf("%f", lat), 1)
	url = strings.Replace(url, "{lon}", fmt.Sprintf("%f", lon), 1)
	url = strings.Replace(url, "{lang}", lang, 1)

	logrus.Debugln("[WeatherClient] accessing ", url)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

// ErrNoFreeConnection is returned when the maximum number of parallel connections is reached
var ErrNoFreeConnection = fmt.Errorf("heating up, try again later")
