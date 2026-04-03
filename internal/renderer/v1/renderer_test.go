package v1

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
)

var update = flag.Bool("update", false, "update golden files")

// testCase represents a single renderer test case.
type testCase struct {
	Name   string       `json:"name"`
	Query  domain.Query `json:"query"`
	Golden string       `json:"golden"`
}

func TestV1Renderer_Render(t *testing.T) {
	// Load raw weather data (required by the renderer)
	weatherRaw, err := loadWeatherRaw("testdata/weather.json")
	if err != nil {
		t.Fatalf("failed to load weather data: %v", err)
	}

	// Load test cases
	cases, err := loadTestCases("testdata/testcases.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	renderer := NewV1Renderer()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Attach weather data to the query
			tc.Query.Weather = weatherRaw

			// Ensure Options is initialized (prevents nil pointer dereference)
			if tc.Query.Options == nil {
				tc.Query.Options = &options.Options{
					Lang: "en", // default fallback
				}
			}

			// Call renderer (only accepts domain.Query)
			output, err := renderer.Render(tc.Query)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			got := output.Content

			// Load expected golden output
			goldenPath := filepath.Join("testdata", tc.Golden+".txt")
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
			}

			// Update golden files when running with -update flag
			if *update {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
				t.Logf("updated golden file: %s", goldenPath)
				return
			}

			// Compare actual vs expected
			if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
				t.Errorf("output mismatch for test case %q\n\n"+
					"--- GOT ---\n%s\n\n"+
					"--- WANT ---\n%s",
					tc.Name, string(got), string(want))

				// Save actual output for easy debugging
				actualPath := goldenPath + ".actual"
				if writeErr := os.WriteFile(actualPath, got, 0o644); writeErr != nil {
					t.Logf("warning: could not write .actual file: %v", writeErr)
				} else {
					t.Logf("actual output saved to: %s", actualPath)
				}
			}
		})
	}
}

// ===================================================================
// Helpers
// ===================================================================

// loadWeatherRaw loads the weather JSON and converts it to *domain.WeatherRaw
func loadWeatherRaw(path string) (*domain.WeatherRaw, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	wr := domain.WeatherRaw(b)
	return &wr, nil
}

// loadTestCases loads the array of test cases from JSON
func loadTestCases(path string) ([]testCase, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cases []testCase
	if err := json.Unmarshal(b, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}