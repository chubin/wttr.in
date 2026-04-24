// Package termutil provides utilities for working with terminal ANSI output,
// including truecolor (24-bit) to 256-color conversion.
package termutil

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gookit/color"
)

// TruecolorTo256 converts all 24-bit truecolor ANSI sequences ([38;2;R;G;Bm and [48;2;R;G;Bm)
// into the closest 256-color equivalents ([38;5;NNNm and [48;5;NNNm).
func TruecolorTo256(input []byte) []byte {
	if len(input) == 0 {
		return input
	}

	var result bytes.Buffer
	i := 0

	for i < len(input) {
		// Look for CSI sequence: ESC [
		if input[i] == '\x1b' && i+1 < len(input) && input[i+1] == '[' {
			start := i
			i += 2 // skip \x1b[

			// Find end of sequence (marked by 'm')
			end := bytes.IndexByte(input[i:], 'm')
			if end == -1 {
				// Incomplete sequence, write as-is
				result.Write(input[start:])
				break
			}
			end += i

			seq := input[start : end+1] // full sequence \x1b[...m

			converted := convertSequence(seq)
			result.Write(converted)
			i = end + 1
			continue
		}

		result.WriteByte(input[i])
		i++
	}

	return result.Bytes()
}

// convertSequence converts a single ANSI sequence if it contains truecolor.
func convertSequence(seq []byte) []byte {
	s := string(seq)

	// Foreground truecolor: 38;2;R;G;B
	if idx := strings.Index(s, "38;2;"); idx != -1 {
		var r, g, b uint8
		n, _ := fmt.Sscanf(s[idx+5:], "%d;%d;%d", &r, &g, &b)
		if n == 3 {
			idx256 := color.RgbTo256(r, g, b)
			return []byte(fmt.Sprintf("\x1b[38;5;%dm", idx256))
		}
	}

	// Background truecolor: 48;2;R;G;B
	if idx := strings.Index(s, "48;2;"); idx != -1 {
		var r, g, b uint8
		n, _ := fmt.Sscanf(s[idx+5:], "%d;%d;%d", &r, &g, &b)
		if n == 3 {
			idx256 := color.RgbTo256(r, g, b)
			return []byte(fmt.Sprintf("\x1b[48;5;%dm", idx256))
		}
	}

	// Return original sequence if no truecolor found
	return seq
}

func RemoveANSI(text []byte) []byte {
	if len(text) == 0 {
		return text
	}

	// Create a new slice to store the result
	result := make([]byte, 0, len(text))
	i := 0

	for i < len(text) {
		if text[i] == 0x1B { // ESC character
			// Skip until we find a letter (a-z or A-Z) which usually terminates the sequence
			i++
			for i < len(text) && !((text[i] >= 'A' && text[i] <= 'Z') || (text[i] >= 'a' && text[i] <= 'z')) {
				i++
			}
			i++ // Skip the terminating letter
		} else {
			// Add non-ANSI character to result
			result = append(result, text[i])
			i++
		}
	}

	return result
}
