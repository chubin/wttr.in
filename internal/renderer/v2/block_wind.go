// internal/renderer/v2/block_wind.go
package v2

import (
	"fmt"
	"math"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

// Wind direction constants from lib/constants.py
var windDirectionSymbols = []string{
	"↓", "↙", "←", "↖", "↑", "↗", "→", "↘",
}

var windDirectionSymbolsWI = []string{
	"", "", "", "", "", "", "", "",
}

// Color thresholds exactly as in Python draw_wind
var windColorThresholds = []struct {
	speed int
	color int
}{
	{3, 241},
	{6, 242},
	{9, 243},
	{12, 246},
	{15, 250},
	{19, 253},
	{23, 214},
	{27, 208},
	{31, 202},
	{-1, 196},
}

func colorCodeForWindSpeed(speed int) string {
	for _, t := range windColorThresholds {
		if speed <= t.speed {
			return fmt.Sprintf("38;5;%d", t.color)
		}
	}
	return "38;5;196"
}

// drawWind renders the two-line wind block (exactly like Python)
func drawWind(dirs []int, speeds []float64, opts *options.Options) string {
	if len(dirs) == 0 || len(speeds) == 0 {
		return "\n\n"
	}

	// Choose symbol set based on view (matches Python logic)
	useWI := opts.View == "v2n" || opts.View == "v2d"
	var symbols []string
	if useWI {
		symbols = windDirectionSymbolsWI
	} else {
		symbols = windDirectionSymbols
	}

	var dirLine, speedLine strings.Builder

	for i := range dirs {
		degree := dirs[i]
		speed := int(speeds[i])

		// Direction symbol
		var symbol string
		if degree != 0 {
			// Exact same formula as Python: (degree + 22.5) % 360 / 45
			angle := float64(degree) + 22.5
			mod := math.Mod(angle, 360)
			idx := int(mod / 45) % 8
			symbol = symbols[idx]
		} else {
			symbol = " "
		}

		color := colorCodeForWindSpeed(speed)

		// Top line: direction with spaces
		dirLine.WriteString(fmt.Sprintf(" %s ", colorize(symbol, color)))

		// Bottom line: speed with padding (exact Python logic)
		speedStr := fmt.Sprintf("%d", speed)
		if speed < 10 {
			speedStr = " " + speedStr + " "
		} else if speed < 100 {
			speedStr = " " + speedStr
		}
		speedLine.WriteString(colorize(speedStr, color))
	}

	return dirLine.String() + "\n" + speedLine.String() + "\n"
}