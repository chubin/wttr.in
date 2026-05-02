package v2

import (
	"strconv"

	"github.com/chubin/wttr.in/internal/domain"
)

// Helper extractors (clean and reusable)
func extractAllHourlyFloat(w domain.Weather, getter func(domain.Hourly) string) []float64 {
	var out []float64
	for _, day := range w.Weather {
		for _, h := range day.Hourly {
			if v, err := strconv.ParseFloat(getter(h), 64); err == nil {
				out = append(out, v)
			} else {
				out = append(out, 0)
			}
		}
	}
	return out
}

func extractAllHourlyInt(w domain.Weather, getter func(domain.Hourly) string) []int {
	var out []int
	for _, day := range w.Weather {
		for _, h := range day.Hourly {
			if v, err := strconv.Atoi(getter(h)); err == nil {
				out = append(out, v)
			} else {
				out = append(out, 0)
			}
		}
	}
	return out
}

// Simple min/max helpers
func slicesMax(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func slicesMin(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
