package v2

import (
	"github.com/kerrigan29a/drawille-go"
)

// drawTemperatureDiagram draws a clean, compact Braille temperature plot.
func drawTemperatureDiagram(data []float64, heightChars, widthChars int) string {
	heightDots := heightChars * 4
	widthDots := widthChars * 2

	minT := slicesMin(data)
	maxT := slicesMax(data)
	rangeT := maxT - minT
	if rangeT == 0 {
		rangeT = 1
	}

	c := drawille.NewCanvas()
	c.Inverse = false

	// Plot the line using exact width/height
	prevX, prevY := -1, -1
	for i, t := range data {
		x := i * (widthDots - 1) / (len(data) - 1)
		normalized := (t - minT) / rangeT
		y := int(normalized * float64(heightDots-1))
		y = heightDots - 1 - y // higher temp = higher on screen

		c.Set(x, y)
		if prevX >= 0 {
			c.DrawLine(prevX, prevY, x, y)
		}
		prevX, prevY = x, y
	}

	// // Label on its own row (below the plot)
	// label := fmt.Sprintf("min:%.0f°   max:%.0f°", minT, maxT)
	// c.SetText(0, heightDots+1, label)

	return c.Frame(0, 0, widthDots-1, heightDots+8)
}
