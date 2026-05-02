package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
)

// colorize returns the string with ANSI background color (46 = cyan bg)
func colorize(s string, colorCode string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, s)
}

func drawTimeScale(loc *domain.Location) string {
	tz, err := time.LoadLocation(loc.TimeZone)
	if err != nil {
		return ""
	}

	// Build top border line (repeated 3 times)
	var line0 strings.Builder
	part := "─────┴─────"
	for i := 0; i < 3; i++ {
		line0.WriteString(part)
		line0.WriteString("┼")
		line0.WriteString(part)
		line0.WriteString("╂")
	}
	line0.WriteRune('\n')

	// Build bottom time scale line
	scaleSegment := "     6    12    18      " // 24 characters
	var line1 strings.Builder
	for i := 0; i < 3; i++ {
		line1.WriteString(scaleSegment)
	}
	line1.WriteRune('\n')

	// Get current hour (0-23) in the given timezone
	now := time.Now().In(tz)
	y, m, d := now.Date()
	correctMidnight := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	hourNumber := int(now.Sub(correctMidnight).Hours())

	// Safety check
	if hourNumber < 0 || hourNumber > 23 {
		hourNumber = 0
	}

	// Highlight current hour on both lines
	s0 := line0.String()
	s1 := line1.String()

	// Convert strings to rune slices for proper Unicode handling
	s0Runes := []rune(s0)
	s1Runes := []rune(s1)

	// Modify s0 if hourNumber is within bounds
	if hourNumber < len(s0Runes) {
		coloredRune := colorize(string(s0Runes[hourNumber]), "46")
		s0 = string(s0Runes[:hourNumber]) + coloredRune + string(s0Runes[hourNumber+1:])
	}

	// Modify s1 if hourNumber is within bounds
	if hourNumber < len(s1Runes) {
		coloredRune := colorize(string(s1Runes[hourNumber]), "46")
		s1 = string(s1Runes[:hourNumber]) + coloredRune + string(s1Runes[hourNumber+1:])
	}

	return s0 + s1
}
