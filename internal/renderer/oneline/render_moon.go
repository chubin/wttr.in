package oneline

import (
	"fmt"
	"math"
	"time"

	"github.com/kixorz/suncalc"
)

// moonAgeDays returns the approximate age of the moon in days
// (0 ≈ new moon, ~14.765 ≈ full moon, ~29.53 ≈ next new moon)
func moonAgeDays(t time.Time) float64 {
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

// RenderMoonPhaseEmoji returns a Unicode moon phase emoji
// Matches common 8-phase style used in wttr.in and many terminal tools
func RenderMoonPhaseEmoji(ctx *RenderContext) string {
	illum := suncalc.GetMoonIllumination(ctx.Now)

	// illum.Phase is 0.0 (new) → 1.0 (next new moon)
	// Map to 8 phases (0=new, 1=wx crescent, ..., 4=full, ..., 7=wn crescent)
	// +0.5 before floor gives nice rounding/centering
	phaseIndex := int(math.Floor(illum.Phase*8+0.5)) % 8

	phases := [8]string{
		0: "🌑", // New Moon
		1: "🌒", // Waxing Crescent
		2: "🌓", // First Quarter
		3: "🌔", // Waxing Gibbous
		4: "🌕", // Full Moon
		5: "🌖", // Waning Gibbous
		6: "🌗", // Last Quarter
		7: "🌘", // Waning Crescent
	}

	return phases[phaseIndex]
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
