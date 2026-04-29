// internal/renderer/v2/block_astronomical.go
package v2

import (
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
)

// drawAstronomical — exact port of the original Python draw_astronomical function
func drawAstronomical(loc *domain.Location, opts *options.Options) string {
	if loc == nil || (loc.Latitude == 0 && loc.Longitude == 0) {
		return ""
	}

	var timeline strings.Builder
	var moonLine strings.Builder

	// City's timezone for correct day start
	locTZ := time.UTC
	if loc.TimeZone != "" {
		if tz, err := time.LoadLocation(loc.TimeZone); err == nil {
			locTZ = tz
		}
	}

	dayStart := time.Now().In(locTZ).Truncate(24 * time.Hour)

	for interval := 0; interval < 72; interval++ {
		current := dayStart.Add(time.Duration(interval) * time.Hour)

		ctx := &oneline.RenderContext{
			Location: loc,
			Options:  opts,
			Now:      current,
		}

		// Use the existing exported Render* functions from oneline package
		dawnStr := oneline.RenderDawn(ctx)
		sunriseStr := oneline.RenderSunrise(ctx)
		sunsetStr := oneline.RenderSunset(ctx)
		duskStr := oneline.RenderDusk(ctx)

		dawnT := parseTimeStr(dawnStr, current)
		sunriseT := parseTimeStr(sunriseStr, current)
		sunsetT := parseTimeStr(sunsetStr, current)
		duskT := parseTimeStr(duskStr, current)

		// === Exact same logic as Python ===
		char := "."
		if current.Before(dawnT) || current.After(duskT) {
			char = " "
		} else if (dawnT.Before(current) || dawnT.Equal(current)) && current.Before(sunriseT) ||
			(sunsetT.Before(current) || sunsetT.Equal(current)) && current.Before(duskT) {
			char = "─"
		} else if (sunriseT.Before(current) || sunriseT.Equal(current)) && current.Before(sunsetT) {
			char = "━"
		}

		timeline.WriteString(char)

		// === Moon phase (exact spacing from Python) ===
		if interval == 0 || interval == 23 || interval == 47 || interval == 69 {
			moonLine.WriteString(oneline.RenderMoonPhaseEmoji(ctx))
		} else if interval%3 == 0 {
			if interval != 24 && interval != 28 {
				moonLine.WriteString("   ")
			} else {
				moonLine.WriteString(" ")
			}
		}
	}

	return moonLine.String() + "\n" + timeline.String() + "\n\n"
}

// parseTimeStr converts "HH:MM:SS" (or "??:??:??") into a time.Time on the same day as `ref`.
// Falls back to `ref` on error — matches Python’s try/except behavior.
func parseTimeStr(s string, ref time.Time) time.Time {
	if s == "" || s == "??:??:??" || s == "??:??:???" {
		return ref
	}

	// Parse time of day
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return ref
	}

	// Build full timestamp on the reference date + location
	return time.Date(
		ref.Year(), ref.Month(), ref.Day(),
		t.Hour(), t.Minute(), t.Second(), 0,
		ref.Location(),
	)
}
