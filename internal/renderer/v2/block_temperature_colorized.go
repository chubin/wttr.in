package v2

import (
	"fmt"
	"math"
	"strings"

	"github.com/kerrigan29a/drawille-go"
	"github.com/rcarmo/go-te/pkg/te"

	"github.com/chubin/wttr.in/internal/renderer/teansi"
)

const (
	FixedMinTemp = -15.0 // °C (very cold)
	FixedMaxTemp = 35.0  // °C (very hot)
)

// colorForTemp returns a temperature-based color (cold blue → hot red).
// You can tweak the gradient here.
func colorForTemp(t, minT, maxT float64) te.Color {
	if maxT == minT {
		return teansi.TrueColor(100, 200, 255) // neutral cyan
	}
	normalized := (t - minT) / (maxT - minT)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	var r, g, b float64

	switch {
	case normalized <= 0.25: // deep blue → cyan
		r = 0
		g = normalized * 4 * 255
		b = 255
	case normalized <= 0.5: // cyan → green
		r = 0
		g = 255
		b = 255 - (normalized-0.25)*4*255
	case normalized <= 0.75: // green → yellow
		r = (normalized - 0.5) * 4 * 255
		g = 255
		b = 0
	default: // yellow → red
		r = 255
		g = 255 - (normalized-0.75)*4*255
		b = 0
	}

	return teansi.TrueColor(uint8(r), uint8(g), uint8(b))
}

// DrawColoredTemperatureDiagram returns the Braille temperature plot
// with color gradient + high/low temperature labels.
func DrawColoredTemperatureDiagram(data []float64, heightChars, widthChars int) string {
	if len(data) == 0 {
		return "No temperature data"
	}

	heightDots := heightChars * 4
	widthDots := widthChars * 2

	minT := slicesMin(data)
	maxT := slicesMax(data)
	rangeT := maxT - minT
	if rangeT == 0 {
		rangeT = 1
	}

	// 1. Build Braille plot
	c := drawille.NewCanvas()
	c.Inverse = false
	prevX, prevY := -1, -1
	for i, t := range data {
		x := i * (widthDots - 1) / (len(data) - 1)
		normalized := (t - minT) / rangeT
		y := int(normalized * float64(heightDots-1))
		y = heightDots - 1 - y // invert Y for canvas

		c.Set(x, y)
		if prevX >= 0 {
			c.DrawLine(prevX, prevY, x, y)
		}
		prevX, prevY = x, y
	}

	plotStr := c.Frame(0, 0, widthDots-1, heightDots-1)

	// 2. Load into go-te screen (extra row for labels)
	screen := te.NewScreen(widthChars, heightChars+2) // +2 for labels
	stream := te.NewStream(screen, false)
	normalized := strings.ReplaceAll(plotStr, "\n", "\r\n")
	stream.Feed(normalized)

	// 3. Color the temperature line
	for col := 0; col < widthChars; col++ {
		xDot := col * 2
		idx := float64(xDot) / float64(widthDots-1) * float64(len(data)-1)
		low := int(math.Floor(idx))
		high := low + 1
		if high >= len(data) {
			high = len(data) - 1
		}
		frac := idx - float64(low)
		t := data[low]*(1-frac) + data[high]*frac

		color := colorForTemp(t, FixedMinTemp, FixedMaxTemp)

		for row := 0; row < len(screen.Buffer); row++ {
			if row >= len(screen.Buffer) || col >= len(screen.Buffer[row]) {
				continue
			}
			cell := &screen.Buffer[row][col]
			if cell.Data != " " && cell.Data != "" {
				teansi.SetCellColor(screen, row, col, color, te.Color{Mode: te.ColorDefault})
			}
		}
	}

	// ==================== ADD LABELS ====================

	// High temperature label (top-left)
	highLabel := fmt.Sprintf("▲ %.0f°", maxT)
	teansi.WriteText(screen, 0, 0, highLabel,
		teansi.TrueColor(255, 100, 100), // warm red
		te.Color{Mode: te.ColorDefault})

	// Low temperature label (bottom-left, just below the plot)
	lowLabel := fmt.Sprintf("▼ %.0f°", minT)
	lowRow := heightChars - 1
	teansi.WriteText(screen, lowRow, 0, lowLabel,
		teansi.TrueColor(100, 180, 255), // cold blue
		te.Color{Mode: te.ColorDefault})

	// // Optional: Current/last temperature on the right
	// if len(data) > 0 {
	// 	lastT := data[len(data)-1]
	// 	lastLabel := fmt.Sprintf("%.0f°", lastT)
	// 	// Place near the right edge
	// 	rightCol := widthChars - len(lastLabel) - 1
	// 	if rightCol < 0 {
	// 		rightCol = 0
	// 	}
	// 	teansi.WriteText(screen, 0, rightCol, lastLabel,
	// 		colorForTemp(lastT, FixedMinTemp, FixedMaxTemp),
	// 		te.Color{Mode: te.ColorDefault})
	// }

	return teansi.ToANSI(screen)
}
