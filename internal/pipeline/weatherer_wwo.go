package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WeatherClient is a struct that implements the Weatherer interface to fetch weather data
// from a specified HTTP endpoint using latitude, longitude, and language parameters.
type WeatherClient struct {
	baseURL string
}

// NewWeatherClient creates a new instance of WeatherClient with the provided base URL.
// The baseURL should contain placeholders for latitude, longitude, and language (e.g., "lat={lat}&lon={lon}&lang={lang}").
// Parameters will be replaced during the request.
func NewWeatherClient(baseURL string) *WeatherClient {
	return &WeatherClient{baseURL: baseURL}
}

// GetWeather fetches weather data from the specified URL by replacing the latitude, longitude,
// and language placeholders in the baseURL. It performs an HTTP GET request and returns the response body
// as a byte slice, or an error if the request fails.
func (wc *WeatherClient) GetWeather(lat, lon float64, lang string) ([]byte, error) {
	// Replace placeholders in the baseURL with actual values
	url := strings.Replace(wc.baseURL, "{lat}", fmt.Sprintf("%f", lat), 1)
	url = strings.Replace(url, "{lon}", fmt.Sprintf("%f", lon), 1)
	url = strings.Replace(url, "{lang}", lang, 1)

	// Create an HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging indefinitely
	}

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close() // Correct way to close the response body

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}
