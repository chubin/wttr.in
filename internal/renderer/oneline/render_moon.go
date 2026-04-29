package oneline

import (
	"fmt"
	"math"
	"time"

	"github.com/kixorz/suncalc"
)

// MoonAgeDays returns the approximate age of the moon in days
// (0 ≈ new moon, ~14.765 ≈ full moon, ~29.53 ≈ next new moon)
func MoonAgeDays(t time.Time) float64 {
	// Reference new moon: 2000-01-06 18:14 UTC (Unix timestamp ≈ 947116440)
	const referenceNewMoon = 947116440.0
	const synodicMonth = 29.530588853 // mean synodic month in days

	secondsSinceRef := float64(t.Unix()) - referenceNewMoon
	daysSinceRef := secondsSinceRef / 86400.0

	age := math.Mod(daysSinceRef, synodicMonth)
	if age < 0 {
		age += synodicMonth
	}

	return age
}

// RenderMoonPhaseEmoji — improved to better match visual expectation + Python spirit
func RenderMoonPhaseEmoji(ctx *RenderContext) string {
	if ctx == nil {
		return "🌕"
	}

	illum := suncalc.GetMoonIllumination(ctx.Now)

	// Best mapping: use illumination directly (most accurate for appearance)
	// +0.5 gives good rounding to nearest visual phase
	idx := int(math.Floor(illum.Phase*8+0.5)) % 8

	phases := [8]string{
		"🌑", // 0 New
		"🌒", // 1 Waxing Crescent
		"🌓", // 2 First Quarter
		"🌔", // 3 Waxing Gibbous
		"🌕", // 4 Full
		"🌖", // 5 Waning Gibbous
		"🌗", // 6 Last Quarter
		"🌘", // 7 Waning Crescent
	}

	return phases[idx]
}

// RenderMoonDay returns approximate lunar day (days since last new moon, 0–29)
// Rounded to nearest integer — suitable for weather/terminal display
func RenderMoonDay(ctx *RenderContext) string {
	illum := suncalc.GetMoonIllumination(ctx.Now)

	// Approximate age from phase value:
	// This is a common approximation; real age = synodic * phase, but phase is already normalized
	ageDays := 29.530588853 * illum.Phase

	// Round to nearest day (0–29 range)
	day := int(math.Round(ageDays))
	if day < 0 {
		day = 0
	}
	if day > 29 {
		day = 29
	}

	return fmt.Sprintf("%d", day)
}
