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

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/translate"
)

var update = flag.Bool("update", false, "update golden files")

// testCase represents a single renderer test case.
type testCase struct {
	Name   string       `json:"name"`
	Query  domain.Query `json:"query"`
	Golden string       `json:"golden"`
}

// testData represents the new structure of testcases.json
type testData struct {
	Location  domain.Location `json:"location"`
	TestCases []testCase      `json:"testcases"`
}

func TestV1Renderer_Render(t *testing.T) {
	weatherRaw, err := loadWeatherRaw("testdata/weather.json")
	if err != nil {
		t.Fatalf("failed to load weather data: %v", err)
	}

	td, err := loadTestData("testdata/testcases.json")
	if err != nil {
		t.Fatalf("failed to load test data: %v", err)
	}

	localizer := translate.NewBundle(assets.FS)
	renderer := NewV1Renderer()

	for _, tc := range td.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Use the shared location for every test case
			tc.Query.Location = &td.Location

			// Attach weather data
			tc.Query.Weather = weatherRaw

			// Ensure Options is never nil and set default language
			if tc.Query.Options == nil {
				tc.Query.Options = &options.Options{Lang: "en"}
			}

			output, err := renderer.Render(tc.Query, localizer)
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

				// Show clean line-by-line diff
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
			diff.WriteString(fmt.Sprintf(" GOT : %s$\n", gLine))
			diff.WriteString(fmt.Sprintf(" WANT: %s$\n", wLine))
			diff.WriteString(" ---\n")
		}
	}
	return diff.String()
}

// TestV1Renderer_LocationRTLMark checks that the "Location:" line is
// prefixed with a right-to-left mark (U+200F) for RTL languages, and
// that it stays plain for LTR languages. Without the mark, terminals
// display RTL location names in the wrong visual order (issue #932).
func TestV1Renderer_LocationRTLMark(t *testing.T) {
	weatherRaw, err := loadWeatherRaw("testdata/weather.json")
	if err != nil {
		t.Fatalf("failed to load weather data: %v", err)
	}

	td, err := loadTestData("testdata/testcases.json")
	if err != nil {
		t.Fatalf("failed to load test data: %v", err)
	}

	localizer := translate.NewBundle(assets.FS)
	renderer := NewV1Renderer()

	cases := []struct {
		lang    string
		wantRTL bool
	}{
		{"en", false},
		{"de", false},
		{"fa", true},
		{"ar", true},
		{"he", true},
	}

	for _, tc := range cases {
		t.Run(tc.lang, func(t *testing.T) {
			query := domain.Query{
				Location: &td.Location,
				Weather:  weatherRaw,
				Options:  &options.Options{Lang: tc.lang},
			}

			output, err := renderer.Render(query, localizer)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			content := string(output.Content)

			idx := strings.Index(content, td.Location.FullAddress)
			if idx == -1 {
				t.Fatalf("output does not contain the location's full address:\n%s", content)
			}

			hasRLMBefore := strings.HasSuffix(content[:idx], rlm)

			if hasRLMBefore != tc.wantRTL {
				t.Errorf("lang %q: RLM immediately before location = %v, want %v", tc.lang, hasRLMBefore, tc.wantRTL)
			}
		})
	}
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

// loadTestData loads the new structured testcases.json
func loadTestData(path string) (testData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return testData{}, err
	}

	var td testData
	if err := json.Unmarshal(b, &td); err != nil {
		return testData{}, err
	}

	return td, nil
}
