// internal/oneline/render_solar.go
package oneline

import (
	"fmt"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// sunriseTime returns local sunrise time in "HH:MM" format (24h)
// Uses github.com/nathan-osman/go-sunrise for accurate calculation
func sunriseTime(ctx *renderContext) string {
	if ctx.Location == nil || ctx.Data == nil {
		return "??:??"
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
	return localRise.Format("15:04")
}

func sunsetTime(ctx *renderContext) string {
	if ctx.Location == nil || ctx.Data == nil {
		return "??:??"
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
	return localSet.Format("15:04")
}

// solarNoonTime — approximate as midpoint between sunrise & sunset
func solarNoonTime(ctx *renderContext) string {
	riseStr := sunriseTime(ctx)
	setStr := sunsetTime(ctx)

	if riseStr == "??:??" || setStr == "??:??" {
		return "??:??"
	}

	riseH, riseM, _ := parseHHMM(riseStr)
	setH, setM, _ := parseHHMM(setStr)

	totalMin := (setH*60 + setM) - (riseH*60 + riseM)
	noonMin := (riseH*60 + riseM) + totalMin/2

	noonH := noonMin / 60
	noonM := noonMin % 60

	return fmt.Sprintf("%02d:%02d", noonH, noonM)
}

// Civil dawn/dusk ≈ sunrise/sunset ± ~30–50 min (rough fallback)
// For real civil twilight use github.com/goastro/twilight or similar
func dawnTime(ctx *renderContext) string {
	riseStr := sunriseTime(ctx)
	if riseStr == "??:??" {
		return "??:??"
	}

	h, m, ok := parseHHMM(riseStr)
	if !ok {
		return "??:??"
	}

	// Crude -50 minutes (civil dawn ≈ -6° solar depression)
	m -= 50
	if m < 0 {
		h--
		m += 60
	}
	if h < 0 {
		h += 24
	}

	return fmt.Sprintf("%02d:%02d", h, m)
}

func duskTime(ctx *renderContext) string {
	setStr := sunsetTime(ctx)
	if setStr == "??:??" {
		return "??:??"
	}

	h, m, ok := parseHHMM(setStr)
	if !ok {
		return "??:??"
	}

	// Crude +50 minutes
	m += 50
	if m >= 60 {
		h++
		m -= 60
	}
	h %= 24

	return fmt.Sprintf("%02d:%02d", h, m)
}

// parseHHMM parses "15:04" style string into hours and minutes
func parseHHMM(s string) (h, m int, ok bool) {
	_, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	return h, m, err == nil
}
