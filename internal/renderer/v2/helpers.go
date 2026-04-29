package v2

import (
	"strconv"
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
)

// drawDate - respects NoCaption / NoCity indirectly through frame
func drawDate(loc *domain.Location) string {
	tz, _ := time.LoadLocation(loc.TimeZone)
	now := time.Now().In(tz)

	var b strings.Builder
	for i := 0; i < 3; i++ {
		d := now.AddDate(0, 0, i)
		dateStr := d.Format("Mon 02 Jan")
		padding := (24 - len(dateStr)) / 2
		line := strings.Repeat(" ", padding) + dateStr + strings.Repeat(" ", 24-len(dateStr)-padding)
		b.WriteString(line)
	}
	return b.String()
}

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
