package v2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
)

type mockLocalizer struct{}

func (m *mockLocalizer) Text(lang, key string) string {
	return key
}

func (m *mockLocalizer) Condition(lang string, code int) string {
	return ""
}

func (m *mockLocalizer) ConditionByName(lang, englishName string) string {
	return englishName
}

func (m *mockLocalizer) File(lang, name string) (string, error) {
	return "", nil
}

func buildTestWeatherJSON(precipMM, precipInches string) domain.WeatherRaw {
	hourlyEntries := make([]string, 0, 8)
	for i := 0; i < 8; i++ {
		hourlyEntries = append(hourlyEntries, fmt.Sprintf(`{
			"time": "%d",
			"tempC": "15",
			"tempF": "59",
			"precipMM": "%s",
			"precipInches": "%s",
			"chanceofrain": "65",
			"weatherCode": "116",
			"windspeedKmph": "12",
			"windspeedMiles": "7",
			"winddirDegree": "180"
		}`, i*300, precipMM, precipInches))
	}

	dayJSON := fmt.Sprintf(`{
		"date": "2026-10-05",
		"hourly": [%s],
		"astronomy": [{
			"sunrise": "06:00 AM",
			"sunset": "06:00 PM",
			"moonrise": "08:00 PM",
			"moonset": "07:00 AM"
		}]
	}`, strings.Join(hourlyEntries, ","))

	raw := fmt.Sprintf(`{
		"current_condition": [{
			"temp_C": "15",
			"temp_F": "59",
			"precipMM": "%s",
			"precipInches": "%s",
			"weatherCode": "116",
			"weatherDesc": [{"value": "Partly cloudy"}],
			"humidity": "60",
			"pressure": "1013",
			"windspeedKmph": "12",
			"windspeedMiles": "7"
		}],
		"weather": [%s, %s, %s]
	}`, precipMM, precipInches, dayJSON, dayJSON, dayJSON)

	return domain.WeatherRaw(raw)
}

func TestV2Renderer_RenderPrecipitationMetric(t *testing.T) {
	testCases := []struct {
		name string
		opts *options.Options
	}{
		{
			name: "default options without flags",
			opts: &options.Options{Location: "Berlin", Lang: "en"},
		},
		{
			name: "explicit metric option",
			opts: &options.Options{Location: "Berlin", Lang: "en", UseMetric: true},
		},
		{
			name: "metric option with ms wind",
			opts: &options.Options{Location: "Berlin", Lang: "en", UseMetric: true, UseMsForWind: true},
		},
	}

	weatherRaw := buildTestWeatherJSON("12.7", "0.50")
	renderer := NewV2Renderer()
	localizer := &mockLocalizer{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := domain.Query{
				Location: &domain.Location{Name: "Berlin", TimeZone: "UTC"},
				Options:  tc.opts,
				Weather:  &weatherRaw,
			}

			out, err := renderer.Render(query, localizer)
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}

			content := string(out.Content)
			if !strings.Contains(content, "12.70mm|65%") {
				t.Fatalf("expected output to contain metric precip label %q, but got:\n%s", "12.70mm|65%", content)
			}
			if strings.Contains(content, "0.50in|") {
				t.Fatalf("metric output should not contain imperial precip label '0.50in|'")
			}
		})
	}
}

func TestV2Renderer_RenderPrecipitationImperial(t *testing.T) {
	testCases := []struct {
		name string
		opts *options.Options
	}{
		{
			name: "use_imperial flag",
			opts: &options.Options{Location: "Pittsburgh", Lang: "en", UseImperial: true},
		},
		{
			name: "use_uscs flag",
			opts: &options.Options{Location: "Pittsburgh", Lang: "en", UseUscs: true},
		},
		{
			name: "both imperial and uscs flags",
			opts: &options.Options{Location: "Pittsburgh", Lang: "en", UseImperial: true, UseUscs: true},
		},
	}

	weatherRaw := buildTestWeatherJSON("12.7", "0.50")
	renderer := NewV2Renderer()
	localizer := &mockLocalizer{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := domain.Query{
				Location: &domain.Location{Name: "Pittsburgh", TimeZone: "America/New_York"},
				Options:  tc.opts,
				Weather:  &weatherRaw,
			}

			out, err := renderer.Render(query, localizer)
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}

			content := string(out.Content)
			if !strings.Contains(content, "0.50in|65%") {
				t.Fatalf("expected output to contain imperial precip label %q, but got:\n%s", "0.50in|65%", content)
			}
			if strings.Contains(content, "12.70mm|") {
				t.Fatalf("imperial output should not contain metric precip label '12.70mm|'")
			}
		})
	}
}

func TestV2Renderer_RenderErrors(t *testing.T) {
	renderer := NewV2Renderer()
	localizer := &mockLocalizer{}

	t.Run("nil weather", func(t *testing.T) {
		_, err := renderer.Render(domain.Query{}, localizer)
		if err == nil {
			t.Fatal("expected error for nil weather")
		}
	})

	t.Run("invalid weather JSON", func(t *testing.T) {
		invalidJSON := domain.WeatherRaw("{invalid")
		_, err := renderer.Render(domain.Query{Weather: &invalidJSON}, localizer)
		if err == nil {
			t.Fatal("expected error for malformed weather JSON")
		}
	})

	t.Run("empty forecast days", func(t *testing.T) {
		emptyDays := domain.WeatherRaw(`{"weather": []}`)
		_, err := renderer.Render(domain.Query{Weather: &emptyDays}, localizer)
		if err == nil {
			t.Fatal("expected error for empty forecast days")
		}
	})
}
