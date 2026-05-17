// internal/renderer/v2/renderer.go
package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/localization"
	"github.com/chubin/wttr.in/internal/options"
)

// V2Renderer implements the rich panel view (temperature diagram, rain sparkline,
// weather emojis, colored wind, astronomical data, frame, etc.)
type V2Renderer struct{}

// NewV2Renderer creates a renderer for the v2 rich panel weather view.
func NewV2Renderer() *V2Renderer {
	return &V2Renderer{}
}

type Localize func(text string) string

// Render converts the query's weather JSON into the v2 terminal weather report.
func (r *V2Renderer) Render(q domain.Query, localizer localization.Localizer) (domain.RenderOutput, error) {
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
	l10n := localization.New(localizer, q.Options)

	var buf bytes.Buffer

	// Location not found warning banner
	if loc != nil && loc.LocationNotFound {
		buf.WriteString(drawLocationNotFoundBanner(loc.OriginalLocation, l10n))
	}

	// Date header (3 days)
	buf.WriteString("\n\n")
	buf.WriteString(drawDate(loc, l10n))

	// Temperature diagram
	tempValues := extractAllHourlyFloat(weather, func(h domain.Hourly) string {
		if opts.UseImperial || opts.UseUscs {
			return h.TempF
		}
		return h.TempC
	})
	tempInterp := interpolate(tempValues, 72)
	buf.WriteString("\n")
	buf.WriteString(DrawColoredTemperatureDiagram(tempInterp, 10, 72))

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
	content := addFrame(buf.String(), 72, opts, l10n)

	if !opts.Quiet && !opts.Superquiet && !opts.NoTerminal {
		content += textualInformation(&q, loc, opts, l10n)
	}

	return domain.RenderOutput{
		Content: []byte(content),
	}, nil
}

// ===================================================================
// Remaining drawing functions (not in helpers.go)
// ===================================================================

// drawLocationNotFoundBanner renders a warning banner when the requested
// location could not be found and Oymyakon is used as fallback.
// Mirrors the v1 404 page design with red background text.
func drawLocationNotFoundBanner(originalLocation string, l10n localization.L10n) string {
	const (
		bgRed   = "\033[41m"
		fgWhite = "\033[97m"
		bold    = "\033[1m"
		reset   = "\033[0m"
	)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(bgRed + fgWhite + bold)
	sb.WriteString("╔══════════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║ ")
	msg := l10n.Text("LOCATION_NOT_FOUND")
	if msg == "" || msg == "LOCATION_NOT_FOUND" {
		sb.WriteString("We were unable to find your location")
		if originalLocation != "" {
			sb.WriteString(": " + originalLocation)
		}
	} else {
		sb.WriteString(msg)
		if originalLocation != "" {
			sb.WriteString(": " + originalLocation)
		}
	}
	sb.WriteString("\n")
	sb.WriteString("║ so we have brought you to Oymyakon,\n")
	sb.WriteString("║ one of the coldest permanently inhabited locales on the planet.\n")
	sb.WriteString("╚══════════════════════════════════════════════════════════════════════╝")
	sb.WriteString(reset)
	sb.WriteString("\n")
	return sb.String()
}

func addFrame(content string, width int, opts *options.Options, l10n localization.L10n) string {
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

	title := l10n.Text("CAPTION_WEATHER_REPORT_FOR")
	if opts.Superquiet {
		title = ""
	} else if !opts.Quiet || !opts.NoCity {
		title += " " + opts.Location
	} else if opts.Quiet {
		title = opts.Location
	}

	caption := ""
	if !opts.Superquiet {
		caption = "┤ " + title + " ├"
	}

	frameTop := "┌" + caption + strings.Repeat("─", width-utf8.RuneCountInString(caption)) + "┐\n"

	return frameTop + strings.Join(lines, "\n") + "\n" +
		"└" + strings.Repeat("─", width) + "┘\n"
}
