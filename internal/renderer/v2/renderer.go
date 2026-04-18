// internal/renderer/v2/renderer.go
package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
)

// V2Renderer implements the rich panel view (temperature diagram, rain sparkline,
// weather emojis, colored wind, astronomical data, frame, etc.)
type V2Renderer struct{}

func NewV2Renderer() *V2Renderer {
	return &V2Renderer{}
}

func (r *V2Renderer) Render(q domain.Query) (domain.RenderOutput, error) {
	if q.Weather == nil || len(*q.Weather) == 0 {
		return domain.RenderOutput{}, fmt.Errorf("no weather data available")
	}

	var weather domain.Weather
	if err := json.Unmarshal(*q.Weather, &weather); err != nil {
		return domain.RenderOutput{}, fmt.Errorf("failed to parse weather JSON: %w", err)
	}

	if len(weather.Weather) == 0 {
		return domain.RenderOutput{}, fmt.Errorf("no forecast days available")
	}

	opts := q.Options
	loc := q.Location

	var buf bytes.Buffer

	// Date header (3 days)
	buf.WriteString(drawDate(loc))
	buf.WriteString("\n\n")

	// Temperature diagram
	tempValues := extractAllHourlyFloat(weather, func(h domain.Hourly) string {
		if opts.UseImperial || opts.UseUscs {
			return h.TempF
		}
		return h.TempC
	})
	tempInterp := interpolate(tempValues, 72)
	buf.WriteString(drawTemperatureDiagram(tempInterp, 10, 72))
	buf.WriteString("\n")

	// Time scale
	buf.WriteString(drawTimeScale(loc))
	buf.WriteString("\n")

	// Rain sparkline
	precipValues := extractAllHourlyFloat(weather, func(h domain.Hourly) string { return h.PrecipMM })
	chanceValues := extractAllHourlyFloat(weather, func(h domain.Hourly) string { return h.ChanceOfRain })
	precipInterp := interpolate(precipValues, 72)
	chanceInterp := interpolate(chanceValues, 72)
	buf.WriteString(drawRainSpark(precipInterp, chanceInterp, 5, 72))
	buf.WriteString("\n")

	// Weather emojis
	codes := extractAllHourlyInt(weather, func(h domain.Hourly) string { return h.WeatherCode })
	buf.WriteString(drawWeatherEmoji(codes, opts))
	buf.WriteString("\n")

	// Wind
	windDirs := extractAllHourlyInt(weather, func(h domain.Hourly) string { return h.WinddirDegree })
	windSpeeds := extractAllHourlyFloat(weather, func(h domain.Hourly) string {
		if opts.UseImperial || opts.UseUscs {
			return h.WindspeedMiles
		}
		return h.WindspeedKmph
	})
	buf.WriteString(drawWind(windDirs, windSpeeds, opts))
	buf.WriteString("\n")

	// Astronomical
	buf.WriteString(drawAstronomical(loc))
	buf.WriteString("\n\n")

	// Frame + optional textual information
	content := addFrame(buf.String(), 72, opts)

	if !opts.Quiet && !opts.Superquiet && !opts.NoTerminal {
		content += textualInformation(weather, loc, opts)
	}

	return domain.RenderOutput{
		Content: []byte(content),
	}, nil
}

// ===================================================================
// Remaining drawing functions (not in helpers.go)
// ===================================================================

func drawTimeScale(loc *domain.Location) string {
	return "   6  12  18    6  12  18    6  12  18  \n"
}

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

func drawWind(dirs []int, speeds []float64, opts *options.Options) string {
	var dirLine, speedLine strings.Builder

	for i, deg := range dirs {
		dirSymbol := getWindDirection(deg)
		color := getWindColor(speeds[i])

		dirLine.WriteString(fmt.Sprintf(" %s ", colorize(dirSymbol, color)))

		spd := int(speeds[i])
		speedStr := fmt.Sprintf("%d", spd)
		if spd < 10 {
			speedStr = " " + speedStr + " "
		}
		speedLine.WriteString(colorize(speedStr, color))
	}

	dirLine.WriteRune('\n')
	speedLine.WriteRune('\n')
	return dirLine.String() + speedLine.String()
}

func getWindDirection(deg int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	return dirs[(deg+22)%360/45]
}

func getWindColor(speed float64) string {
	switch {
	case speed < 5:
		return "38;5;242"
	case speed < 15:
		return "38;5;250"
	case speed < 25:
		return "38;5;226"
	default:
		return "38;5;196"
	}
}

func colorize(text, colorCode string) string {
	if colorCode == "" {
		return text
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
}

func drawAstronomical(loc *domain.Location) string {
	// Placeholder - can be expanded with real sunrise/sunset calculation later
	return "─ Sunrise ───── Noon ────── Sunset ───── Dusk ──\n" +
		"   06:12       13:05        20:58        22:10   \n"
}

func addFrame(content string, width int, opts *options.Options) string {
	if opts.NoCaption {
		return content
	}

	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	for i := range lines {
		spacesNumber := width - len(lines[i])
		if spacesNumber < 0 {
			spacesNumber = 0
		}
		lines[i] = "│" + lines[i] + strings.Repeat(" ", spacesNumber) + "│"
	}

	title := " Weather Report "
	if opts.Superquiet {
		title = ""
	} else if opts.Quiet || opts.NoCity {
		title = " " + opts.Location + " "
	}

	caption := "┤" + title + "├"
	frameTop := "┌" + caption + strings.Repeat("─", width-len(caption)) + "┐\n"

	return frameTop + strings.Join(lines, "\n") + "\n" +
		"└" + strings.Repeat("─", width) + "┘\n"
}

func textualInformation(weather domain.Weather, loc *domain.Location, opts *options.Options) string {
	if len(weather.CurrentCondition) == 0 {
		return ""
	}
	curr := weather.CurrentCondition[0]

	desc := ""
	if len(curr.WeatherDesc) > 0 {
		desc = curr.WeatherDesc[0].Value
	}

	return fmt.Sprintf(
		"\nWeather: %s, %s°C\nLocation: %s, %s\nTimezone: %s\n",
		desc, curr.TempC, loc.Name, loc.Country, loc.TimeZone,
	)
}

// interpolate - linear interpolation (used by temperature and rain)
func interpolate(data []float64, targetWidth int) []float64 {
	if len(data) == 0 {
		return make([]float64, targetWidth)
	}
	result := make([]float64, targetWidth)
	n := len(data) - 1
	for i := range result {
		x := float64(i) / float64(targetWidth-1) * float64(n)
		low := int(math.Floor(x))
		high := low + 1
		if high >= len(data) {
			high = len(data) - 1
		}
		frac := x - float64(low)
		result[i] = data[low]*(1-frac) + data[high]*frac
	}
	return result
}
