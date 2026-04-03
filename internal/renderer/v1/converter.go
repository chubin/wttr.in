package v1

import (
	"strconv"

	"github.com/chubin/wttr.in/internal/domain"
)

// ConvertWeather converts domain.Weather to v1.resp
func ConvertWeather(w domain.Weather) resp {
	r := resp{}

	// Current condition
	if len(w.CurrentCondition) > 0 {
		cc := w.CurrentCondition[0]
		r.Data.Cur = []cond{convertCondFromCurrent(cc)}
	}

	// Request
	if len(w.Request) > 0 {
		r.Data.Req = []loc{
			{
				Query: w.Request[0].Query,
				Type:  w.Request[0].Type,
			},
		}
	}

	// Weather days
	for _, wd := range w.Weather {
		day := weather{
			Date:     wd.Date,
			MaxtempC: parseInt(wd.MaxTempC),
			MintempC: parseInt(wd.MinTempC),
		}

		// Astronomy
		if len(wd.Astronomy) > 0 {
			a := wd.Astronomy[0]
			day.Astronomy = []astro{
				{
					Moonrise: a.Moonrise,
					Moonset:  a.Moonset,
					Sunrise:  a.Sunrise,
					Sunset:   a.Sunset,
				},
			}
		}

		// Hourly
		for _, h := range wd.Hourly {
			day.Hourly = append(day.Hourly, convertCondFromHourly(h))
		}

		r.Data.Weather = append(r.Data.Weather, day)
	}

	return r
}

// Helper: convert current_condition
func convertCondFromCurrent(cc domain.CurrentCondition) cond {
	return cond{
		FeelsLikeC:     parseInt(cc.FeelsLikeC),
		PrecipMM:       parseFloat32(cc.PrecipMM),
		TempC2:         parseInt(cc.TempC), // note: temp_C in current
		VisibleDistKM:  parseInt(cc.Visibility),
		WeatherCode:    parseInt(cc.WeatherCode),
		WindspeedKmph:  parseInt(cc.WindspeedKmph),
		Winddir16Point: cc.Winddir16Point,
		WeatherDesc:    convertValueItems(cc.WeatherDesc),
		// ChanceOfRain not present in current_condition → leave empty
	}
}

// Helper: convert hourly
func convertCondFromHourly(h domain.Hourly) cond {
	return cond{
		ChanceOfRain:   h.ChanceOfRain,
		FeelsLikeC:     parseInt(h.FeelsLikeC),
		PrecipMM:       parseFloat32(h.PrecipMM),
		TempC:          parseInt(h.TempC),
		Time:           parseInt(h.Time),
		VisibleDistKM:  parseInt(h.Visibility),
		WeatherCode:    parseInt(h.WeatherCode),
		WindGustKmph:   parseInt(h.WindGustKmph),
		WindspeedKmph:  parseInt(h.WindspeedKmph),
		Winddir16Point: h.Winddir16Point,
		WeatherDesc:    convertValueItems(h.WeatherDesc),
	}
}

func convertValueItems(items []domain.ValueItem) []struct{ Value string } {
	result := make([]struct{ Value string }, len(items))
	for i, item := range items {
		result[i].Value = item.Value
	}
	return result
}

// Safe parsing helpers
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}
