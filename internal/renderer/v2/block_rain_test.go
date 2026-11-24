package v2

import (
	"strings"
	"testing"

	"github.com/chubin/wttr.in/internal/options"
)

func TestDrawRainSparkMetricUnits(t *testing.T) {
	testCases := []struct {
		name         string
		opts         *options.Options
		expectedUnit string
	}{
		{
			name:         "nil options defaults to metric mm",
			opts:         nil,
			expectedUnit: "mm|",
		},
		{
			name:         "empty options defaults to metric mm",
			opts:         &options.Options{},
			expectedUnit: "mm|",
		},
		{
			name:         "explicit use_metric true",
			opts:         &options.Options{UseMetric: true},
			expectedUnit: "mm|",
		},
		{
			name:         "use_metric with use_ms_for_wind",
			opts:         &options.Options{UseMetric: true, UseMsForWind: true},
			expectedUnit: "mm|",
		},
	}

	data := []float64{0.0, 2.5, 12.7, 5.0, 0.0}
	chance := []float64{0, 20, 80, 40, 0}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := drawRainSpark(data, chance, 5, 20, tc.opts)
			if !strings.Contains(output, "12.70mm|80%") {
				t.Fatalf("expected output to contain %q, got:\n%s", "12.70mm|80%", output)
			}
			if !strings.Contains(output, tc.expectedUnit) {
				t.Fatalf("expected output to contain unit %q, got:\n%s", tc.expectedUnit, output)
			}
		})
	}
}

func TestDrawRainSparkImperialUnits(t *testing.T) {
	testCases := []struct {
		name         string
		opts         *options.Options
		expectedUnit string
	}{
		{
			name:         "use_imperial true",
			opts:         &options.Options{UseImperial: true},
			expectedUnit: "in|",
		},
		{
			name:         "use_uscs true",
			opts:         &options.Options{UseUscs: true},
			expectedUnit: "in|",
		},
		{
			name:         "both use_imperial and use_uscs true",
			opts:         &options.Options{UseImperial: true, UseUscs: true},
			expectedUnit: "in|",
		},
	}

	data := []float64{0.0, 0.1, 0.5, 0.2, 0.0}
	chance := []float64{0, 20, 80, 40, 0}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := drawRainSpark(data, chance, 5, 20, tc.opts)
			if !strings.Contains(output, "0.50in|80%") {
				t.Fatalf("expected output to contain %q, got:\n%s", "0.50in|80%", output)
			}
			if !strings.Contains(output, tc.expectedUnit) {
				t.Fatalf("expected output to contain unit %q, got:\n%s", tc.expectedUnit, output)
			}
		})
	}
}

func TestDrawRainSparkEdgeCases(t *testing.T) {
	t.Run("empty data slice", func(t *testing.T) {
		out := drawRainSpark([]float64{}, []float64{}, 5, 20, nil)
		expected := strings.Repeat(" ", 20) + "\n"
		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("zero width", func(t *testing.T) {
		out := drawRainSpark([]float64{1.0, 2.0}, []float64{10, 20}, 5, 0, nil)
		if out != "\n" {
			t.Fatalf("expected %q, got %q", "\n", out)
		}
	})

	t.Run("all zero data values", func(t *testing.T) {
		out := drawRainSpark([]float64{0.0, 0.0, 0.0}, []float64{10, 20, 30}, 5, 15, nil)
		expected := strings.Repeat(" ", 15) + "\n"
		if out != expected {
			t.Fatalf("expected %q, got %q", expected, out)
		}
	})

	t.Run("chanceData shorter than maxCol", func(t *testing.T) {
		// maxValue at index 3, chanceData only has 2 elements
		data := []float64{0.0, 1.0, 2.0, 8.5}
		chance := []float64{10, 20}
		out := drawRainSpark(data, chance, 5, 20, &options.Options{UseMetric: true})
		if !strings.Contains(out, "8.50mm|0%") {
			t.Fatalf("expected %q with fallback 0%% chance, got:\n%s", "8.50mm|0%", out)
		}
	})

	t.Run("label alignment near right edge", func(t *testing.T) {
		// maxValue at the very end of line
		data := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 5.0}
		chance := []float64{0, 0, 0, 0, 0, 90}
		out := drawRainSpark(data, chance, 3, 12, &options.Options{UseMetric: true})
		if !strings.Contains(out, "5.00mm|90%") {
			t.Fatalf("expected right aligned label in output, got:\n%s", out)
		}
	})
}
