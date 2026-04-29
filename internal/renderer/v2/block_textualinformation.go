// internal/renderer/v2/block_textualinformation.go
package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
)

const (
	dim   = "\033[2m"
	reset = "\033[0m"
)

// dimLabel returns a dimmed label like "Weather:" (matches Python v2 style)
func dimLabel(label string) string {
	return dim + label + reset
}

// renderLocalTimeOnly returns only HH:MM:SS without timezone offset
func renderLocalTimeOnly(ctx *oneline.RenderContext) string {
	if ctx == nil || ctx.Location == nil {
		return time.Now().Format("15:04:05")
	}

	loc := time.UTC
	if ctx.Location.TimeZone != "" {
		tz, err := time.LoadLocation(ctx.Location.TimeZone)
		if err == nil {
			loc = tz
		}
	}

	localTime := ctx.Now.In(loc)
	return localTime.Format("15:04:05")
}

// textualInformation returns the rich bottom metadata block (matches original v2 exactly)
func textualInformation(q *domain.Query, loc *domain.Location, opts *options.Options) string {
	if q.Weather == nil || len(*q.Weather) == 0 {
		return simpleTextualFallback(loc)
	}

	// Use oneline parser + renderers for consistency
	data, err := oneline.ParseCurrentCondition(*q.Weather)
	if err != nil {
		return simpleTextualFallback(loc)
	}

	ctx := &oneline.RenderContext{
		Data:     data,
		DataRaw:  nil,
		Options:  opts,
		Location: loc,
		Now:      time.Now(),
	}

	var b strings.Builder

	// === Main weather line ===
	b.WriteString(dimLabel("Weather:") + " ")
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
	b.WriteRune('\n')

	// === Timezone ===
	b.WriteString(dimLabel("Timezone:") + " ")
	b.WriteString(loc.TimeZone)
	b.WriteRune('\n')

	// === Full Astronomy Line ===
	b.WriteString(dimLabel("  Now:   ") + " ")
	b.WriteString(renderLocalTimeOnly(ctx)) // ← Only local time, no +0200
	b.WriteString(" " + dim + "|" + reset + " ")
	b.WriteString(dimLabel("Dawn:  ") + " ")
	b.WriteString(oneline.RenderDawn(ctx))
	b.WriteString(" " + dim + "|" + reset + " ")
	b.WriteString(dimLabel("Sunrise:") + " ")
	b.WriteString(oneline.RenderSunrise(ctx))
	b.WriteRune('\n')

	b.WriteString(dimLabel("  Zenith:") + " ")
	b.WriteString(oneline.RenderSolarNoon(ctx))
	b.WriteString(" " + dim + "|" + reset + " ")
	b.WriteString(dimLabel("Sunset:") + " ")
	b.WriteString(oneline.RenderSunset(ctx))
	b.WriteString(" " + dim + "|" + reset + " ")
	b.WriteString(dimLabel("Dusk:") + " ")
	b.WriteString(oneline.RenderDusk(ctx))
	b.WriteRune('\n')

	// === Location ===
	b.WriteString(dimLabel("Location:") + " ")
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
	b.WriteRune('\n')

	return b.String()
}

// simpleTextualFallback is used when parsing fails
func simpleTextualFallback(loc *domain.Location) string {
	return dimLabel("Weather:") + " ???\n" +
		dimLabel("Timezone:") + " " + loc.TimeZone + "\n" +
		dimLabel("Location:") + " " + loc.Name + "\n"
}
