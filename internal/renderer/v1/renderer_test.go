package v1

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	weatherRaw, err := loadWeatherRaw("testdata/weather.json")
	if err != nil {
		t.Fatalf("failed to load weather data: %v", err)
	}

	cases, err := loadTestCases("testdata/testcases.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	renderer := NewV1Renderer()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Query.Weather = weatherRaw

			if tc.Query.Options == nil {
				tc.Query.Options = &options.Options{Lang: "en"}
			}

			output, err := renderer.Render(tc.Query)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			got := output.Content

			goldenPath := filepath.Join("testdata", tc.Golden+".txt")
			wantBytes, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
			}

			if *update {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
				t.Logf("updated golden file: %s", goldenPath)
				return
			}

			want := bytes.TrimSpace(wantBytes)
			gotTrim := bytes.TrimSpace(got)

			if !bytes.Equal(gotTrim, want) {
				t.Errorf("output mismatch for test case %q", tc.Name)

				// Show colored original output
				t.Logf("\n--- GOT (colored) ---\n%s", string(got))
				t.Logf("\n--- WANT (colored) ---\n%s", string(wantBytes))

				// Show clean line-by-line diff (ANSI stripped)
				diff := diffLines(string(gotTrim), string(want))
				if diff != "" {
					t.Logf("\n--- LINE-BY-LINE DIFF ---\n%s", diff)
				}

				// Save actual output for easy comparison
				actualPath := goldenPath + ".actual"
				if err := os.WriteFile(actualPath, got, 0o644); err == nil {
					t.Logf("actual output saved to: %s", actualPath)
				}
			}
		})
	}
}

// diffLines returns a human-readable diff of two multi-line strings
func diffLines(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")

	var diff strings.Builder
	maxLen := len(g)
	if len(w) > maxLen {
		maxLen = len(w)
	}

	for i := 0; i < maxLen; i++ {
		var gLine, wLine string
		if i < len(g) {
			gLine = g[i]
		}
		if i < len(w) {
			wLine = w[i]
		}

		if gLine != wLine {
			diff.WriteString(fmt.Sprintf("Line %3d:\n", i+1))
			diff.WriteString(fmt.Sprintf("  GOT : %s$\n", gLine))
			diff.WriteString(fmt.Sprintf("  WANT: %s$\n", wLine))
			diff.WriteString("  ---\n")
		}
	}
	return diff.String()
}

// ===================================================================
// Helpers
// ===================================================================

func loadWeatherRaw(path string) (*domain.WeatherRaw, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	wr := domain.WeatherRaw(b)
	return &wr, nil
}

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
