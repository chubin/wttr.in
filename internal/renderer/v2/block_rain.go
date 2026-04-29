// internal/renderer/v2/block_rain.go
package v2

import (
	"fmt"
	"math"
	"strings"
)

// sparkBars as []rune — correct Unicode handling
var sparkBars = []rune(" _▁▂▃▄▅▇█")

// drawRainSpark — now matches Python v2 label position (TOP of the sparkline)
func drawRainSpark(data, chanceData []float64, height, width int) string {
	if len(data) == 0 || width == 0 {
		return strings.Repeat(" ", width) + "\n"
	}

	maxValue := slicesMax(data)
	if maxValue == 0 {
		return strings.Repeat(" ", width) + "\n"
	}

	var output strings.Builder

	// === Build the spark graph (top to bottom) ===
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			rowFromBottom := height - i - 1
			character := box(rowFromBottom, data[j], maxValue, height)

			if data[j] != 0 && j < len(chanceData) && chanceData[j] > 0 {
				chanceOfRain := math.Min(1.0, chanceData[j]/100.0*2.0)
				colorIndex := int(5 * chanceOfRain)
				colorCode := 16 + colorIndex

				output.WriteString(fmt.Sprintf("\033[38;5;%dm%s\033[0m", colorCode, character))
			} else {
				output.WriteString(character)
			}
		}
		output.WriteRune('\n')
	}

	// === Max value label (placed at the TOP, like original Python) ===
	maxCol := 0
	for j, v := range data {
		if v == maxValue {
			maxCol = j
			break
		}
	}

	label := fmt.Sprintf("%3.2fmm|%d%%", maxValue, int(chanceData[maxCol]))
	coloredLabel := fmt.Sprintf("\033[38;5;33m%s\033[0m", label)

	// Python-style alignment
	labelLen := len(label)
	var maxLine string
	if labelLen/2 < maxCol && labelLen/2+maxCol < width {
		spaces := strings.Repeat(" ", maxCol-labelLen/2)
		maxLine = spaces + coloredLabel + strings.Repeat(" ", width-(maxCol-labelLen/2)-labelLen)
	} else if labelLen/2+maxCol >= width {
		maxLine = strings.Repeat(" ", width-labelLen) + coloredLabel
	} else {
		left := max(0, (width-labelLen)/2)
		maxLine = strings.Repeat(" ", left) + coloredLabel + strings.Repeat(" ", width-left-labelLen)
	}

	// Prepend label + blank line (matches original Python behavior)
	return "\n" + maxLine + "\n" + output.String()
}

// box — exact port of Python _box
func box(row int, value, maxValue float64, height int) string {
	if maxValue == 0 {
		return " "
	}

	rowHeight := maxValue / float64(height)

	if rowHeight*float64(row) >= value {
		return " "
	}
	if rowHeight*float64(row+1) <= value {
		return string(sparkBars[len(sparkBars)-1])
	}

	filled := (value - rowHeight*float64(row)) / rowHeight
	idx := int(filled * float64(len(sparkBars)))
	if idx >= len(sparkBars) {
		idx = len(sparkBars) - 1
	}
	return string(sparkBars[idx])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
