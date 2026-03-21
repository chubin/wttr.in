package oneline

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

// ──────────────────────────────────────────────────────────────────────────────
// All render functions — same signature
// ──────────────────────────────────────────────────────────────────────────────

func renderConditionFullName(ctx *renderContext) string {
	if ctx.Options.Lang == "en" {
		return ctx.Data.WeatherDesc
	}

	query := fmt.Sprintf(".data.current_condition[0].lang_%s[0].value", ctx.Options.Lang)

	val, err := GetStringValue(ctx.DataRaw, query)
	if err != nil {
		log.Println("renderConditionFullName: ", err)
		return ctx.Data.WeatherDesc
	}

	return val
}

func renderConditionCode(ctx *renderContext) string {
	return ctx.Data.ConditionCode
}

func renderTemperature(ctx *renderContext) string {
	var val float64
	var unit string
	if ctx.Options.UseImperial {
		val, unit = ctx.Data.TempF, "°F"
	} else {
		val, unit = ctx.Data.TempC, "°C"
	}
	sign := ""
	if val >= 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%.0f%s", sign, math.Round(val), unit)
}

func renderFeelsLike(ctx *renderContext) string {
	var val float64
	var unit string
	if ctx.Options.UseImperial {
		val, unit = ctx.Data.FeelsLikeF, "°F"
	} else {
		val, unit = ctx.Data.FeelsLikeC, "°C"
	}
	sign := ""
	if val >= 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%.0f%s", sign, math.Round(val), unit)
}

func renderWind(ctx *renderContext) string {
	// dir := windDirSymbol(ctx.Data.WindDirDegree, ctx.Options.View)
	dir := ""

	var speed float64
	var unit string
	switch {
	case ctx.Options.UseMsForWind:
		speed = ctx.Data.WindKmph / 3.6
		unit = "m/s"
	case ctx.Options.UseImperial:
		speed = ctx.Data.WindMiles
		unit = "mph"
	default:
		speed = ctx.Data.WindKmph
		unit = "km/h"
	}

	return fmt.Sprintf("%s%.0f%s", dir, math.Round(speed), unit)
}

func renderHumidity(ctx *renderContext) string {
	return fmt.Sprintf("%d%%", ctx.Data.Humidity)
}

func renderPrecipitation(ctx *renderContext) string {
	return fmt.Sprintf("%.1fmm", ctx.Data.PrecipMM)
}

func renderPrecipChance(ctx *renderContext) string {
	return fmt.Sprintf("%d%%", ctx.Data.ChanceOfRain)
}

func renderPressure(ctx *renderContext) string {
	return fmt.Sprintf("%dhPa", ctx.Data.PressureHpa)
}

func renderUVIndex(ctx *renderContext) string {
	return strconv.Itoa(ctx.Data.UVIndex)
}

func renderLocation(ctx *renderContext) string {
	if ctx.Location != nil {
		return ctx.Location.Name
	}
	if ctx.Data.LocationName != "" {
		return ctx.Data.LocationName
	}
	return "?"
}

func renderSunrise(ctx *renderContext) string {
	// implement with geo + date
	return sunriseTime(ctx) // stub or real impl
}

func renderSunset(ctx *renderContext) string {
	return sunsetTime(ctx)
}

func renderDawn(ctx *renderContext) string {
	return dawnTime(ctx)
}

func renderDusk(ctx *renderContext) string {
	return duskTime(ctx)
}

func renderSolarNoon(ctx *renderContext) string {
	return solarNoonTime(ctx)
}

func renderLocalTime(ctx *renderContext) string {
	// Start with UTC as default
	loc := time.UTC
	if ctx.Location.TimeZone != "" {
		tz, err := time.LoadLocation(ctx.Location.TimeZone)
		if err == nil {
			loc = tz
		}
	}

	// Convert the time to the specified timezone
	localTime := ctx.Now.In(loc)
	return localTime.Format("15:04:05-0700")
}

func renderTimezone(ctx *renderContext) string {
	if ctx.Location != nil && ctx.Location.TimeZone != "" {
		return ctx.Location.TimeZone
	}
	return "UTC"
}

// renderConditionPlain implements the %x placeholder
// Returns a very simple plain-text (ASCII-like) representation of the weather condition
func renderConditionPlain(ctx *renderContext) string {
	if ctx == nil || ctx.Data == nil || ctx.Data.ConditionCode == "" {
		return "?"
	}

	// Step 1: Get symbolic name from code (e.g. "116" → "PartlyCloudy")
	symbolicName, ok := WWOCodeToName[ctx.Data.ConditionCode]
	if !ok {
		return "?"
	}

	// Step 2: Look up the plain ASCII-style symbol
	plainSymbol, exists := WeatherSymbolPlain[symbolicName]
	if !exists {
		return "?"
	}

	return plainSymbol
}
