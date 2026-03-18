package oneline

import (
	"fmt"
	"math"
	"time"
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

// renderMoonPhaseEmoji returns a Unicode moon phase emoji for the given time
// (8 main phases — matches most weather services and wttr.in style)
func renderMoonPhaseEmoji(ctx *renderContext) string {
	age := moonAgeDays(ctx.Now)

	// Map age (0–29.53) → 0–7 phase index
	// We use simple linear mapping with slight bias toward better visual distinction
	phase := int(math.Floor((age / 29.530588853 * 8) + 0.5))

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

	return phases[phase%8]
}

// renderMoonDay returns days since last new moon (0–29)
// Rounded to nearest integer — common in weather displays
func renderMoonDay(ctx *renderContext) string {
	age := moonAgeDays(ctx.Now)
	day := int(math.Round(age))
	// Sometimes people cap at 28, but 0–29 is more accurate for most purposes
	return fmt.Sprintf("%d", day)
}
