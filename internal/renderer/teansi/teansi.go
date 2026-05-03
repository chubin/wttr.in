// Package teansi adds additional features to go-te, as for example
// ANSI serialization and easy color manipulation.
package teansi

import (
	"fmt"
	"strings"

	"github.com/rcarmo/go-te/pkg/te"
)

// ToANSI converts a go-te Screen back into a colored ANSI string
// that can be printed directly to a terminal.
func ToANSI(s *te.Screen) string {
	if s == nil || len(s.Buffer) == 0 {
		return ""
	}

	var buf strings.Builder
	lastFg := te.Color{Mode: te.ColorDefault}
	lastBg := te.Color{Mode: te.ColorDefault}

	for row := range s.Buffer {
		for col := range s.Buffer[row] {
			cell := s.Buffer[row][col]

			// Only emit color codes when they change
			if !colorsEqual(cell.Attr.Fg, lastFg) {
				buf.WriteString(colorToANSI(cell.Attr.Fg, true))
				lastFg = cell.Attr.Fg
			}
			if !colorsEqual(cell.Attr.Bg, lastBg) {
				buf.WriteString(colorToANSI(cell.Attr.Bg, false))
				lastBg = cell.Attr.Bg
			}

			buf.WriteString(cell.Data)
		}
		buf.WriteString("\033[0m\n") // reset + newline
		lastFg = te.Color{Mode: te.ColorDefault}
		lastBg = te.Color{Mode: te.ColorDefault}
	}

	return buf.String()
}

// SetCellColor sets foreground/background color for a specific cell.
func SetCellColor(s *te.Screen, row, col int, fg, bg te.Color) {
	if row < 0 || row >= len(s.Buffer) || col < 0 || col >= len(s.Buffer[row]) {
		return
	}
	cell := &s.Buffer[row][col]
	if fg.Mode != te.ColorDefault {
		cell.Attr.Fg = fg
	}
	if bg.Mode != te.ColorDefault {
		cell.Attr.Bg = bg
	}
}

// SetLineColor is useful for your Braille plot: colors an entire line (or part of it).
func SetLineColor(s *te.Screen, row int, fg te.Color) {
	if row < 0 || row >= len(s.Buffer) {
		return
	}
	for col := range s.Buffer[row] {
		if s.Buffer[row][col].Data != " " { // only color non-blank cells
			s.Buffer[row][col].Attr.Fg = fg
		}
	}
}

// Helper: create a true-color (RGB)
func TrueColor(r, g, b uint8) te.Color {
	return te.Color{
		Mode: te.ColorTrueColor,
		Name: fmt.Sprintf("#%02x%02x%02x", r, g, b),
	}
}

// Helper: create 256-color
func ANSI256Color(index uint8) te.Color {
	return te.Color{
		Mode:  te.ColorANSI256,
		Index: index,
	}
}

// colorToANSI converts internal Color to ANSI escape sequence.
func colorToANSI(c te.Color, isFg bool) string {
	prefix := "38" // foreground
	if !isFg {
		prefix = "48" // background
	}

	switch c.Mode {
	case te.ColorDefault:
		return fmt.Sprintf("\033[%dm", 39) // reset fg or 49 for bg
	case te.ColorANSI16:
		code := c.Index
		if !isFg {
			code += 10
		}
		return fmt.Sprintf("\033[%dm", code)
	case te.ColorANSI256:
		return fmt.Sprintf("\033[%s;5;%dm", prefix, c.Index)
	case te.ColorTrueColor:
		// c.Name is expected to be "#rrggbb"
		if len(c.Name) == 7 && c.Name[0] == '#' {
			r := hexToUint8(c.Name[1:3])
			g := hexToUint8(c.Name[3:5])
			b := hexToUint8(c.Name[5:7])
			return fmt.Sprintf("\033[%s;2;%d;%d;%dm", prefix, r, g, b)
		}
	}
	return ""
}

func colorsEqual(a, b te.Color) bool {
	return a.Mode == b.Mode && a.Index == b.Index && a.Name == b.Name
}

func hexToUint8(s string) uint8 {
	var v uint8
	fmt.Sscanf(s, "%02x", &v)
	return v
}

// WriteText writes a string at (row, col) with optional foreground/background color.
// Stops at screen edges. Supports \n.
func WriteText(s *te.Screen, row, col int, text string, fg, bg te.Color) {
	if s == nil || len(s.Buffer) == 0 {
		return
	}

	r, c := row, col

	for _, ch := range text {
		if ch == '\n' {
			r++
			c = col
			continue
		}

		if r < 0 || r >= len(s.Buffer) || c < 0 || c >= len(s.Buffer[r]) {
			continue // skip characters outside bounds
		}

		cell := &s.Buffer[r][c]
		cell.Data = string(ch)
		SetCellColor(s, r, c, fg, bg)

		c++
	}
}
