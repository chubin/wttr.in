package ansitopng

import (
	"image/color"

	"github.com/rcarmo/go-te/pkg/te"
)

// colorFromANSI converts go-te color (from Attr.Fg / Attr.Bg) to real color.Color
func colorFromANSI(c te.Color, inverse bool) color.Color {
	if c.Mode == te.ColorDefault || c.Name == "default" {
		if inverse {
			return color.Black
		}
		return color.RGBA{211, 211, 211, 255} // lightgray
	}

	// Named color
	switch c.Name {
	case "black":
		return color.Black
	case "red":
		return color.RGBA{205, 49, 49, 255}
	case "green":
		return color.RGBA{0, 128, 0, 255}
	case "yellow":
		return color.RGBA{205, 205, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 205, 255}
	case "magenta":
		return color.RGBA{205, 0, 205, 255}
	case "cyan":
		return color.RGBA{0, 205, 205, 255}
	case "white":
		return color.RGBA{229, 229, 229, 255}
	}

	// 256-color mode (most important for wttr.in)
	if c.Mode == te.ColorANSI256 {
		return ansi256ToColor(int(c.Index))
	}

	return color.RGBA{211, 211, 211, 255}
}

// ansi256ToColor maps common 256-color values used by wttr.in
func ansi256ToColor(n int) color.Color {
	switch {
	case n == 226 || n == 227 || n == 190 || n == 220: // yellow / sun
		return color.RGBA{255, 215, 0, 255}
	case n >= 46 && n <= 82: // green shades
		return color.RGBA{0, 255, 0, 255}
	case n >= 118 && n <= 154: // bright green / cyan
		return color.RGBA{0, 255, 100, 255}
	case n >= 196 && n <= 208: // red / orange
		return color.RGBA{255, 60, 60, 255}
	case n >= 33 && n <= 39: // blue
		return color.RGBA{100, 100, 255, 255}
	}
	return color.RGBA{200, 200, 200, 255}
}
