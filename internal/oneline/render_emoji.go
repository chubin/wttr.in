package oneline

import "strconv"

// renderConditionEmoji returns a weather condition emoji for the %c placeholder
// Supports day/night variants via ctx.Options.View ("v2d", "v2n", or default)
func renderConditionEmoji(ctx *renderContext) string {
	codeStr := ctx.Data.ConditionCode
	if codeStr == "" {
		return "❓"
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return "❓"
	}

	view := ""
	if ctx.Options != nil {
		view = ctx.Options.View // "v2d" = day icons, "v2n" = night icons
	}

	// ────────────────────────────────────────────────────────────────
	// Core mapping (WorldWeatherOnline codes → emoji)
	// Most common codes — extend as needed
	// ────────────────────────────────────────────────────────────────
	switch code {
	// Clear / Sunny
	case 113:
		if view == "v2n" {
			return "🌙" // clear night
		}
		return "☀️" // sunny day

	// Partly cloudy
	case 116:
		if view == "v2n" {
			return "🌙☁️" // or just "☁️" — many use same
		}
		return "⛅"

	// Cloudy
	case 119, 122:
		return "☁️"

	// Overcast
	case 123:
		return "🌫️" // or "☁️☁️"

	// Mist / Fog
	case 248, 260:
		return "🌫️"

	// Patchy rain / light rain
	case 176, 293, 296, 353:
		if view == "v2n" {
			return "🌙🌧️"
		}
		return "🌦️"

	// Rain / moderate/heavy rain
	case 302, 308, 359, 356:
		return "🌧️"

	// Heavy rain showers / torrential
	case 305, 371:
		return "⛈️" // or "🌧️🌧️"

	// Drizzle
	case 266, 281, 284:
		return "🌧️" // light version

	// Snow / snow showers
	case 323, 326, 329, 332, 335, 377:
		if view == "v2n" {
			return "🌙❄️"
		}
		return "❄️"

	// Sleet / mixed rain+snow
	case 317, 320, 350:
		return "🌨️" // or "❄️🌧️"

	// Thunderstorm / thundery showers
	case 200, 386, 389, 392, 395:
		return "⛈️"

	// default fallback
	default:
		return "🌫️" // unknown → neutral symbol
	}
}
