package v2

import "strings"

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
