// internal/renderer/v2/renderer.go
package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

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
	buf.WriteString("\n\n")
	buf.WriteString(drawDate(loc))

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
	buf.WriteString(drawAstronomical(loc, opts))
	buf.WriteString("\n\n")

	// Frame + optional textual information
	content := addFrame(buf.String(), 72, opts)

	if !opts.Quiet && !opts.Superquiet && !opts.NoTerminal {
		content += textualInformation(&q, loc, opts)
	}

	return domain.RenderOutput{
		Content: []byte(content),
	}, nil
}

// ===================================================================
// Remaining drawing functions (not in helpers.go)
// ===================================================================

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

	title := "  Weather report for: "
	if opts.Superquiet {
		title = ""
	} else if !opts.Quiet || !opts.NoCity {
		title += opts.Location + "  "
	} else if opts.Quiet {
		title = opts.Location + "  "
	}

	caption := "┤" + title + "├"
	frameTop := "┌" + caption + strings.Repeat("─", width-utf8.RuneCountInString(caption)) + "┐\n"

	return frameTop + strings.Join(lines, "\n") + "\n" +
		"└" + strings.Repeat("─", width) + "┘\n"
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
