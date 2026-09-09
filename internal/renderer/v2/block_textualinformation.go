// internal/renderer/v2/block_textualinformation.go
package v2

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/localization"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
)

const (
	dim   = "\033[2m"
	reset = "\033[0m"

	// rlm is the Unicode right-to-left mark (U+200F). See its use in
	// buildLocationString and simpleTextualFallback below.
	rlm = "‏"
)

func dimLabel(label string) string {
	return dim + label + reset
}

// displayWidth returns the number of columns a string occupies on a terminal.
// It correctly handles emojis, wide characters, combining characters, etc.
func displayWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size

		switch {
		case r == '\t':
			width += 8 // tab width
		case r < 0x20: // control characters
			continue
		default:
			// Simple but effective heuristic used by many terminal tools
			if r >= 0x1100 && (r <= 0x115f || // Hangul Jamo
				r == 0x2329 || r == 0x232a ||
				(r >= 0x2e80 && r <= 0x4dbf) || // CJK
				(r >= 0x4e00 && r <= 0x9fff) ||
				(r >= 0xac00 && r <= 0xd7a3) ||
				(r >= 0xf900 && r <= 0xfaff) ||
				(r >= 0xfe10 && r <= 0xfe19) ||
				(r >= 0xfe30 && r <= 0xfe6f) ||
				(r >= 0xff00 && r <= 0xff60) ||
				(r >= 0xffe0 && r <= 0xffe6)) {
				width += 2
			} else if r >= 0x1f300 && r <= 0x1f9ff { // Emojis
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// renderNowWithOffset returns current local time with timezone offset
func renderNowWithOffset(ctx *oneline.RenderContext) string {
	if ctx == nil || ctx.Location == nil {
		return time.Now().Format("15:04:05-0700")
	}

	loc := time.UTC
	if ctx.Location.TimeZone != "" {
		tz, err := time.LoadLocation(ctx.Location.TimeZone)
		if err == nil {
			loc = tz
		}
	}
	return ctx.Now.In(loc).Format("15:04:05-0700")
}

type field struct {
	label string
	value string
}

type section []field

// textualInformation returns the rich bottom metadata block (data-first)
func textualInformation(q *domain.Query, loc *domain.Location, opts *options.Options, l10n localization.L10n) string {
	if q.Weather == nil || len(*q.Weather) == 0 {
		return simpleTextualFallback(loc, l10n.IsRTL())
	}

	data, err := oneline.ParseCurrentCondition(*q.Weather)
	if err != nil {
		return simpleTextualFallback(loc, l10n.IsRTL())
	}

	ctx := &oneline.RenderContext{
		Data:     data,
		Options:  opts,
		Location: loc,
		Now:      time.Now(),
	}

	// === CLEAN DATA DEFINITION ===
	sections := []section{
		{{label: l10n.Text("WEATHER"), value: buildMainWeatherLine(ctx)}},
		{{label: l10n.Text("TIMEZONE"), value: loc.TimeZone}},
		{
			{label: l10n.Text("NOW"), value: renderNowWithOffset(ctx)},
			{label: l10n.Text("DAWN"), value: oneline.RenderDawn(ctx)},
			{label: l10n.Text("SUNRISE"), value: oneline.RenderSunrise(ctx)},
		},
		{
			{label: l10n.Text("ZENITH"), value: oneline.RenderSolarNoon(ctx)},
			{label: l10n.Text("SUNSET"), value: oneline.RenderSunset(ctx)},
			{label: l10n.Text("DUSK"), value: oneline.RenderDusk(ctx)},
		},
		{{label: l10n.Text("LOCATION"), value: buildLocationString(loc, l10n.IsRTL())}},
	}

	return renderSections(sections)
}

// renderSections with proper terminal display width
func renderSections(sections []section) string {
	var b strings.Builder

	// === Step 1: Calculate GLOBAL max display width per column ===
	maxLabel := make([]int, 0)
	maxValue := make([]int, 0)

	for _, sec := range sections {
		if len(sec) <= 1 {
			continue
		}

		// Extend slices if needed
		for len(maxLabel) < len(sec) {
			maxLabel = append(maxLabel, 0)
			maxValue = append(maxValue, 0)
		}

		for j, f := range sec {
			if w := displayWidth(f.label); w > maxLabel[j] {
				maxLabel[j] = w
			}
			if w := displayWidth(f.value); w > maxValue[j] {
				maxValue[j] = w
			}
		}
	}

	// === Step 2: Render with alignment ===
	for _, sec := range sections {
		if len(sec) == 0 {
			continue
		}

		if len(sec) == 1 {
			f := sec[0]
			b.WriteString(dimLabel(f.label+": ") + f.value + "\n")
			continue
		}

		// Multi-field section
		for j, f := range sec {
			labelWidth := displayWidth(f.label)
			padding := maxLabel[j] - labelWidth

			// Label (right-aligned) + colon (part of dimmed label)
			paddedLabel := f.label + strings.Repeat(" ", padding) + ":"
			b.WriteString(dimLabel(paddedLabel) + " ")

			// Value (left-aligned)
			valueWidth := displayWidth(f.value)
			valuePadding := maxValue[j] - valueWidth
			paddedValue := f.value + strings.Repeat(" ", valuePadding)
			b.WriteString(paddedValue)

			if j < len(sec)-1 {
				b.WriteString(" " + dim + "|" + reset + " ")
			}
		}
		b.WriteRune('\n')
	}

	return b.String()
}

// buildMainWeatherLine constructs the rich first line
func buildMainWeatherLine(ctx *oneline.RenderContext) string {
	var b strings.Builder
	b.WriteString(oneline.RenderConditionEmoji(ctx))
	b.WriteString(oneline.RenderConditionFullName(ctx))
	b.WriteString(", ")
	b.WriteString(oneline.RenderTemperature(ctx))
	b.WriteString(", ")
	b.WriteString(oneline.RenderHumidity(ctx))
	b.WriteString(", ")
	b.WriteString(oneline.RenderWind(ctx))
	b.WriteString(", ")
	b.WriteString(oneline.RenderPressure(ctx))
	return b.String()
}

// buildLocationString builds the location line
func buildLocationString(loc *domain.Location, isRTL bool) string {
	var b strings.Builder
	if isRTL {
		// Right-to-left mark (U+200F): without it, terminals display
		// an RTL location name in the wrong visual order (issue #932).
		b.WriteString(rlm)
	}
	if loc.FullAddress != "" {
		b.WriteString(loc.FullAddress)
	} else {
		b.WriteString(loc.Name)
		if loc.Country != "" {
			b.WriteString(", " + loc.Country)
		}
	}

	if loc.Latitude != 0 || loc.Longitude != 0 {
		b.WriteString(fmt.Sprintf(" [%.4f,%.4f]", loc.Latitude, loc.Longitude))
	}
	return b.String()
}

// simpleTextualFallback
func simpleTextualFallback(loc *domain.Location, isRTL bool) string {
	location := loc.Name
	if isRTL {
		location = rlm + location
	}
	return dimLabel("Weather: ") + "???\n" +
		dimLabel("Timezone: ") + loc.TimeZone + "\n" +
		dimLabel("Location: ") + location + "\n"
}