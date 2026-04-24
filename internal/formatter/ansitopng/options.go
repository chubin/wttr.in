package ansitopng

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

// PNGOptions contains all options relevant to ANSI → PNG rendering.
type PNGOptions struct {
	// Background color for the generated PNG (default: black)
	Background color.Color
	// Invert foreground and background colors
	Inverted bool
	// Transparency level for the final PNG (0 = fully transparent, 255 = fully opaque)
	Transparency int
	// Add a frame around the PNG (from ?frame)
	Frame bool
	// Custom width and height for the generated image (0 = auto)
	Width  int
	Height int
	// Add padding around content
	Padding bool
}

// NewPNGOptions returns sensible defaults matching original wttr.in behavior
func NewPNGOptions() PNGOptions {
	return PNGOptions{
		Background:   color.Black,
		Inverted:     false,
		Transparency: 255, // fully opaque by default
		Frame:        false,
		Width:        0,
		Height:       0,
		Padding:      false,
	}
}

// FromDomainOptions converts from the main wttr.in options.Options to PNGOptions
func FromDomainOptions(opts *options.Options) PNGOptions {
	if opts == nil {
		return NewPNGOptions()
	}

	bg := color.Color(color.Black)
	if opts.Background != "" {
		bg = parseColor(opts.Background)
	}

	transparency := 255
	if opts.Transparency != 0 {
		transparency = clamp(opts.Transparency, 0, 255)
	}

	return PNGOptions{
		Background:   bg,
		Inverted:     opts.InvertedColors,
		Transparency: transparency,
		Frame:        opts.Frame,
		Width:        opts.Width,
		Height:       opts.Height,
		Padding:      opts.Padding,
	}
}

// parseColor converts a hex string (RRGGBB) or named color to color.Color
func parseColor(s string) color.Color {
	s = strings.TrimPrefix(strings.ToLower(s), "#")

	switch s {
	case "black":
		return color.Black
	case "white":
		return color.White
	case "lightgray", "lightgrey":
		return color.RGBA{211, 211, 211, 255}
	case "green":
		return color.RGBA{0, 128, 0, 255}
	case "cyan":
		return color.RGBA{0, 255, 255, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "brown":
		return color.RGBA{165, 42, 42, 255}
	}

	// Parse RRGGBB hex
	if len(s) == 6 {
		if r, err := strconv.ParseUint(s[0:2], 16, 8); err == nil {
			if g, err := strconv.ParseUint(s[2:4], 16, 8); err == nil {
				if b, err := strconv.ParseUint(s[4:6], 16, 8); err == nil {
					return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
				}
			}
		}
	}
	return color.Black // safe fallback
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
