package v2

import (
	"fmt"
	"math"
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

// drawTemperatureDiagram - basic block diagram (easy to upgrade later)
func drawTemperatureDiagram(data []float64, height, width int) string {
	if len(data) == 0 {
		return strings.Repeat(" ", width) + "\n"
	}
	maxT := slicesMax(data)
	minT := slicesMin(data)
	if maxT == minT {
		maxT = minT + 10
	}

	var b strings.Builder
	for row := height - 1; row >= 0; row-- {
		level := minT + (maxT-minT)*float64(row+1)/float64(height)
		for _, val := range data {
			if val >= level {
				b.WriteRune('█')
			} else {
				b.WriteRune('░')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

// drawRainSpark - uses color based on chance of rain (same logic as Python)
func drawRainSpark(data, chanceData []float64, height, width int) string {
	bars := " _▁▂▃▄▅▇█"
	var b strings.Builder

	maxVal := slicesMax(data)
	if maxVal == 0 {
		return "\n" + strings.Repeat(" ", width) + "\n"
	}

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			val := data[col]
			ch := getSparkChar(height, height-row-1, val, maxVal, bars)

			colorCode := 20 // default gray
			if col < len(chanceData) && val > 0 {
				chance := math.Min(1.0, chanceData[col]/100.0*2.0)
				colorCode = 16 + int(5*chance)
			}
			b.WriteString(fmt.Sprintf("\033[38;5;%dm%s\033[0m", colorCode, ch))
		}
		b.WriteRune('\n')
	}

	// Max value label
	b.WriteString(fmt.Sprintf("\n%4.1fmm\n", maxVal))
	return b.String()
}

func getSparkChar(height, row int, value, maxValue float64, bars string) string {
	if maxValue == 0 {
		return " "
	}
	rowHeight := maxValue / float64(height)
	if rowHeight*float64(row) >= value {
		return " "
	}
	idx := int(((value - rowHeight*float64(row)) / rowHeight) * float64(len(bars)))
	if idx >= len(bars) {
		idx = len(bars) - 1
	}
	return string(bars[idx])
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
