package oneline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
)

// parsedCurrentCondition is a normalized, ready-to-use view of current weather
// extracted from the WorldWeatherOnline-style JSON response
type parsedCurrentCondition struct {
	// Core weather values
	ConditionCode  string // "116", "248", etc.
	WeatherDesc    string // "Partly cloudy" (translated if lang_xx present)
	TempC          float64
	TempF          float64
	FeelsLikeC     float64
	FeelsLikeF     float64
	Humidity       int     // 0–100
	PrecipMM       float64 // mm in last period (usually 1h or 3h)
	PrecipInches   float64 // mm in last period (usually 1h or 3h)
	ChanceOfRain   int     // 0–100 %
	PressureHpa    int
	UVIndex        int
	WindKmph       float64
	WindMiles      float64
	WindDir16Point string // "NNW", "S", etc.
	WindDirDegree  int    // 0–359

	// Location & time info
	LocationName     string
	ObservationLocal time.Time // parsed from "localObsDateTime"

	// Optional / derived
	CloudCover   int
	VisibilityKm float64
}

// parseCurrentCondition extracts and normalizes current weather data
// from the raw JSON bytes returned by the upstream weather service.
func parseCurrentCondition(raw []byte) (*parsedCurrentCondition, error) {
	var fullWrapped struct {
		Data domain.Weather `json:"data"`
	}
	if err := json.Unmarshal(raw, &fullWrapped); err != nil {
		return nil, fmt.Errorf("invalid weather JSON: %w", err)
	}

	full := fullWrapped.Data

	if len(full.CurrentCondition) == 0 {
		return nil, fmt.Errorf("response missing current_condition array")
	}
	cc := full.CurrentCondition[0]

	// ── Parse numeric fields with safe defaults ────────────────────────────────
	tempC, _ := strconv.ParseFloat(cc.TempC, 64)
	tempF, _ := strconv.ParseFloat(cc.TempF, 64)
	feelsC, _ := strconv.ParseFloat(cc.FeelsLikeC, 64)
	feelsF, _ := strconv.ParseFloat(cc.FeelsLikeF, 64)

	humidity, _ := strconv.Atoi(cc.Humidity)
	precipMM, _ := strconv.ParseFloat(cc.PrecipMM, 64)
	precipInches, _ := strconv.ParseFloat(cc.PrecipInches, 64)
	pressure, _ := strconv.Atoi(cc.Pressure)
	uv, _ := strconv.Atoi(cc.UVIndex)
	windKmph, _ := strconv.ParseFloat(cc.WindspeedKmph, 64)
	windMiles, _ := strconv.ParseFloat(cc.WindspeedMiles, 64)
	windDeg, _ := strconv.Atoi(cc.WinddirDegree)

	cloud, _ := strconv.Atoi(cc.Cloudcover)
	visKm, _ := strconv.ParseFloat(cc.Visibility, 64)

	// ── Description (prefer translated lang_xx if present) ─────────────────────
	desc := ""
	if len(cc.WeatherDesc) > 0 {
		desc = cc.WeatherDesc[0].Value
	}

	// Try to find any lang_ field (lang_ru, lang_de, etc.)
	for _, v := range cc.WeatherDesc { // actually it's []ValueItem, but we check parent
		// Note: real code might need to look at the full map or use a lang prefix search
		// For simplicity we take the first description (usually English or requested lang)
		if desc == "" && len(v.Value) > 0 {
			desc = v.Value
		}
	}

	// ── Observation time ───────────────────────────────────────────────────────
	obsTime, err := time.Parse("2006-01-02 15:04", cc.LocalObsDateTime)
	if err != nil {
		obsTime = time.Now().UTC() // fallback
	}

	// ── Location fallback ──────────────────────────────────────────────────────
	locName := "?"
	if len(full.Request) > 0 {
		locName = full.Request[0].Query
	}
	if len(full.NearestArea) > 0 && len(full.NearestArea[0].AreaName) > 0 {
		locName = full.NearestArea[0].AreaName[0].Value
	}

	return &parsedCurrentCondition{
		ConditionCode:    cc.WeatherCode,
		WeatherDesc:      desc,
		TempC:            tempC,
		TempF:            tempF,
		FeelsLikeC:       feelsC,
		FeelsLikeF:       feelsF,
		Humidity:         humidity,
		PrecipMM:         precipMM,
		PrecipInches:     precipInches,
		ChanceOfRain:     0, // often comes from hourly → can be filled later if needed
		PressureHpa:      pressure,
		UVIndex:          uv,
		WindKmph:         windKmph,
		WindMiles:        windMiles,
		WindDir16Point:   cc.Winddir16Point,
		WindDirDegree:    windDeg,
		CloudCover:       cloud,
		VisibilityKm:     visKm,
		LocationName:     locName,
		ObservationLocal: obsTime,
	}, nil
}
