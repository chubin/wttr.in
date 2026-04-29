package v2

import (
	"fmt"
	"math"
	"strings"
)

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
