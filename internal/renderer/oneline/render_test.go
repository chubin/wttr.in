package oneline

import (
	"math"
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

var dewPointWeatherFixture = domain.WeatherRaw(`{
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
  "request": [{"query": "Thibodaux", "type": "City"}],
  "weather": [{
    "date": "2026-06-04",
    "hourly": [
      {"time": "900", "DewPointC": "11", "DewPointF": "52"},
      {"time": "1200", "DewPointC": "16", "DewPointF": "61"},
      {"time": "1500", "DewPointC": "18", "DewPointF": "64"}
    ]
  }]
}`)

func TestRenderDewPointUsesNearestHourly(t *testing.T) {
	out, err := NewOnelineRenderer().Render(domain.Query{
		Options: &options.Options{Format: "%e"},
		Weather: &dewPointWeatherFixture,
	}, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got, want := string(out.Content), "+16\u00b0C"; got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRenderDewPointUsesImperialUnits(t *testing.T) {
	for name, opts := range map[string]*options.Options{
		"imperial": {Format: "%e", UseImperial: true},
		"uscs":     {Format: "%e", UseUscs: true},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := NewOnelineRenderer().Render(domain.Query{
				Options: opts,
				Weather: &dewPointWeatherFixture,
			}, nil)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if got, want := string(out.Content), "+61\u00b0F"; got != want {
				t.Fatalf("Render() = %q, want %q", got, want)
			}
		})
	}
}

func TestRenderDewPointFallsBackToMagnus(t *testing.T) {
	out, err := NewOnelineRenderer().Render(domain.Query{
		Options: &options.Options{Format: "%e"},
		Weather: &currentWeatherFixture,
	}, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got, want := string(out.Content), "+17\u00b0C"; got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestFormatSignedDegreesZero(t *testing.T) {
	for _, tc := range []struct {
		val  float64
		want string
	}{
		{math.Copysign(0, -1), "0°C"}, // temp_C "-0" from upstream
		{-0.3, "0°C"},
		{0.4, "0°C"},
		{-0.6, "-1°C"},
		{18, "+18°C"},
		{-5, "-5°C"},
	} {
		if got := formatSignedDegrees(tc.val, "°C"); got != tc.want {
			t.Errorf("formatSignedDegrees(%v) = %q, want %q", tc.val, got, tc.want)
		}
	}
}
