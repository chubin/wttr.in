package oneline

import (
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
)

var currentWeatherFixture = domain.WeatherRaw(`{
  "current_condition": [{
    "FeelsLikeC": "10",
    "FeelsLikeF": "50",
    "cloudcover": "25",
    "humidity": "92",
    "localObsDateTime": "2026-06-04 12:00",
    "precipInches": "0.50",
    "precipMM": "12.7",
    "pressure": "1012",
    "temp_C": "18",
    "temp_F": "64",
    "uvIndex": "5",
    "visibility": "10",
    "weatherCode": "116",
    "weatherDesc": [{"value": "Partly cloudy"}],
    "winddir16Point": "S",
    "winddirDegree": "180",
    "windspeedKmph": "10",
    "windspeedMiles": "6"
  }],
  "request": [{"query": "Thibodaux", "type": "City"}]
}`)

func TestRenderPrecipitationUsesInchesForImperialFormats(t *testing.T) {
	for name, opts := range map[string]*options.Options{
		"imperial": {Format: "%p", UseImperial: true},
		"uscs":     {Format: "%p", UseUscs: true},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := NewOnelineRenderer().Render(domain.Query{
				Options: opts,
				Weather: &currentWeatherFixture,
			}, nil)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if got, want := string(out.Content), "0.50 in"; got != want {
				t.Fatalf("Render() = %q, want %q", got, want)
			}
		})
	}
}

func TestRenderPrecipitationKeepsMetricDefault(t *testing.T) {
	out, err := NewOnelineRenderer().Render(domain.Query{
		Options: &options.Options{Format: "%p"},
		Weather: &currentWeatherFixture,
	}, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got, want := string(out.Content), "12.7mm"; got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRenderUSCSUsesImperialWeatherFields(t *testing.T) {
	out, err := NewOnelineRenderer().Render(domain.Query{
		Options: &options.Options{Format: "%t %w %p", UseUscs: true},
		Weather: &currentWeatherFixture,
	}, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got, want := string(out.Content), "+64\u00b0F ↑6mph 0.50 in"; got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}
