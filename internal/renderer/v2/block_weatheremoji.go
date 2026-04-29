package v2

import (
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

func drawWeatherEmoji(codes []int, opts *options.Options) string {
	// Basic weather code to emoji mapping (expand as needed)
	emojiMap := map[int]string{
		113: "☀️", 116: "⛅", 119: "☁️", 122: "☁️",
		176: "🌦️", 200: "⛈️", 227: "❄️", 230: "❄️",
		248: "🌫️", 260: "🌫️",
	}

	var b strings.Builder
	for _, code := range codes {
		emoji := emojiMap[code]
		if emoji == "" {
			emoji = "🌡️"
		}
		if opts.StandardFont {
			emoji = "*"
		}
		b.WriteString(emoji + "  ")
	}
	b.WriteRune('\n')
	return b.String()
}
