// internal/oneline/render_solar.go
package oneline

import (
	"fmt"
	"os"
	"time"

	"github.com/goastro/twilight"
	"github.com/nathan-osman/go-sunrise"
)

// sunriseTime returns local sunrise time in "HH:MM" format (24h)
// Uses github.com/nathan-osman/go-sunrise for accurate calculation
func sunriseTime(ctx *renderContext) string {
	if ctx.Location == nil || ctx.Data == nil {
		return "??:??:??"
	}

	lat := ctx.Location.Latitude
	lon := ctx.Location.Longitude
	date := ctx.Now.Truncate(24 * time.Hour) // midnight of "today" (UTC or local — adjust if needed)

	rise, _ := sunrise.SunriseSunset(
		lat, lon,
		date.Year(), date.Month(), date.Day(),
	)

	// Convert UTC result to local timezone (if available)
	loc := time.UTC
	if ctx.Location.TimeZone != "" {
		tz, err := time.LoadLocation(ctx.Location.TimeZone)
		if err == nil {
			loc = tz
		}
	}

	localRise := rise.In(loc)
	return localRise.Format("15:04:05")
}

func sunsetTime(ctx *renderContext) string {
	if ctx.Location == nil || ctx.Data == nil {
		return "??:??:??"
	}

	lat := ctx.Location.Latitude
	lon := ctx.Location.Longitude
	date := ctx.Now.Truncate(24 * time.Hour)

	_, set := sunrise.SunriseSunset(
		lat, lon,
		date.Year(), date.Month(), date.Day(),
	)

	loc := time.UTC
	if ctx.Location.TimeZone != "" {
		tz, err := time.LoadLocation(ctx.Location.TimeZone)
		if err == nil {
			loc = tz
		}
	}

	localSet := set.In(loc)
	return localSet.Format("15:04:05")
}

// solarNoonTime — approximate as midpoint between sunrise & sunset
func solarNoonTime(ctx *renderContext) string {
	riseStr := sunriseTime(ctx)
	setStr := sunsetTime(ctx)

	if riseStr == "??:??:??" || setStr == "??:??:??" {
		return "??:??:??"
	}

	riseH, riseM, riseS, _ := parseHHMMSS(riseStr)
	setH, setM, setS, _ := parseHHMMSS(setStr)

	totalMin := (setH*60 + setM) - (riseH*60 + riseM)
	noonMin := (riseH*60 + riseM) + totalMin/2

	noonH := noonMin / 60
	noonM := noonMin % 60
	noonS := (riseS - setS + 60) % 60

	return fmt.Sprintf("%02d:%02d:%02d", noonH, noonM, noonS)
}

// getLatLonAndTZ extracts latitude, longitude, and timezone from renderContext
// Returns:
//   - lat, lon: coordinates (0,0 if missing/invalid)
//   - loc: timezone location (UTC fallback if missing)
//   - ok: true only if we have usable coordinates
func getLatLonAndTZ(ctx *renderContext) (lat, lon float64, loc *time.Location, ok bool) {
	if ctx == nil {
		return 0, 0, time.UTC, false
	}

	// 1. Preferred source: Location struct (most reliable when available)
	if ctx.Location != nil {
		lat = ctx.Location.Latitude
		lon = ctx.Location.Longitude

		// Timezone from Location.TimeZone string
		if ctx.Location.TimeZone != "" {
			var err error
			loc, err = time.LoadLocation(ctx.Location.TimeZone)
			if err == nil {
				// success — we have both coords + valid tz
				return lat, lon, loc, lat != 0 || lon != 0
			}
			// log invalid tz name but continue with coords
			fmt.Fprintf(os.Stderr, "WARNING: invalid timezone %q in Location → falling back to UTC\n",
				ctx.Location.TimeZone)
		}
	}

	// 2. Last resort fallback timezone: UTC
	loc = time.UTC

	// Consider coordinates valid only if they're clearly non-zero
	// (you can make this check stricter if needed, e.g. |lat| <= 90, |lon| <= 180)
	ok = lat != 0 || lon != 0

	if !ok {
		fmt.Fprintln(os.Stderr, "WARNING: no usable latitude/longitude in renderContext")
	}

	return lat, lon, loc, ok
}

// ────────────────────────────────────────────────
// Optional: Helper to get "today" at midnight local time
// Useful for twilight calculations
// ────────────────────────────────────────────────
func getLocalMidnight(ctx *renderContext) time.Time {
	_, _, loc, _ := getLatLonAndTZ(ctx) // ignore ok here — UTC is safe fallback
	now := ctx.Now
	if now.IsZero() {
		now = time.Now()
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}

// dawnTime returns civil dawn time in local HH:MM format (or "??:??" on error/failure)
func dawnTime(ctx *renderContext) string {
	lat, lon, loc, ok := getLatLonAndTZ(ctx)
	if !ok {
		return "??:??:??"
	}

	// Use a proper date – here assuming "today" from context or current time
	// In real code, prefer to take date from astronomy/weather data if available

	now := time.Now().In(loc)
	utcDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC) // midnight UTC

	civilRise, _, status := twilight.CivilTwilight(utcDate, lat, lon)

	if status != twilight.SunriseStatusOK {
		return "??:??:??"
	}

	// Then convert the morning one (civilRise) to local
	localDawn := civilRise.In(loc)

	return localDawn.Format("15:04:05") //
}

// duskTime returns civil dusk time in local HH:MM format (or "??:??" on error/failure)
func duskTime(ctx *renderContext) string {
	lat, lon, loc, ok := getLatLonAndTZ(ctx)
	if !ok {
		return "??:??:??"
	}

	now := time.Now().In(loc)
	utcDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	_, civilDusk, status := twilight.CivilTwilight(utcDate, lat, lon)

	if status != twilight.SunriseStatusOK {
		// fmt.Printf("WARNING: no civil dusk on %s at %.2f,%.2f (status: %d)\n",
		// 	date.Format("2006-01-02"), lat, lon, status)
		return "??:??:??"
	}

	localDusk := civilDusk.In(loc)
	return localDusk.Format("15:04:05")
}

// parseHHMM parses "15:04" style string into hours and minutes
func parseHHMMSS(s string) (h, m, ss int, ok bool) {
	_, err := fmt.Sscanf(s, "%d:%d:%d", &h, &m, &ss)
	return h, m, ss, err == nil
}
