package v2

import (
	"fmt"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

func drawWind(dirs []int, speeds []float64, opts *options.Options) string {
	var dirLine, speedLine strings.Builder

	for i, deg := range dirs {
		dirSymbol := getWindDirection(deg)
		color := getWindColor(speeds[i])

		dirLine.WriteString(fmt.Sprintf(" %s ", colorize(dirSymbol, color)))

		spd := int(speeds[i])
		speedStr := fmt.Sprintf("%d", spd)
		if spd < 10 {
			speedStr = " " + speedStr + " "
		}
		speedLine.WriteString(colorize(speedStr, color))
	}

	dirLine.WriteRune('\n')
	speedLine.WriteRune('\n')
	return dirLine.String() + speedLine.String()
}

func getWindDirection(deg int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	return dirs[(deg+22)%360/45]
}

func getWindColor(speed float64) string {
	switch {
	case speed < 5:
		return "38;5;242"
	case speed < 15:
		return "38;5;250"
	case speed < 25:
		return "38;5;226"
	default:
		return "38;5;196"
	}
}
