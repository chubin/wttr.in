// internal/oneline/constants.go
package oneline

import "github.com/enescakir/emoji"

// WWOCodeToName maps WorldWeatherOnline weather codes to internal symbolic names.
// This is the canonical key used by all emoji/icon/width profiles.
var WWOCodeToName = map[string]string{
	"113": "Sunny",
	"116": "PartlyCloudy",
	"119": "Cloudy",
	"122": "VeryCloudy",
	"143": "Fog",
	"176": "LightShowers",
	"179": "LightSleetShowers",
	"182": "LightSleet",
	"185": "LightSleet",
	"200": "ThunderyShowers",
	"227": "LightSnow",
	"230": "HeavySnow",
	"248": "Fog",
	"260": "Fog",
	"263": "LightShowers",
	"266": "LightRain",
	"281": "LightSleet",
	"284": "LightSleet",
	"293": "LightRain",
	"296": "LightRain",
	"299": "HeavyShowers",
	"302": "HeavyRain",
	"305": "HeavyShowers",
	"308": "HeavyRain",
	"311": "LightSleet",
	"314": "LightSleet",
	"317": "LightSleet",
	"320": "LightSnow",
	"323": "LightSnowShowers",
	"326": "LightSnowShowers",
	"329": "HeavySnow",
	"332": "HeavySnow",
	"335": "HeavySnowShowers",
	"338": "HeavySnow",
	"350": "LightSleet",
	"353": "LightShowers",
	"356": "HeavyShowers",
	"359": "HeavyRain",
	"362": "LightSleetShowers",
	"365": "LightSleetShowers",
	"368": "LightSnowShowers",
	"371": "HeavySnowShowers",
	"374": "LightSleetShowers",
	"377": "LightSleet",
	"386": "ThunderyShowers",
	"389": "ThunderyHeavyRain",
	"392": "ThunderySnowShowers",
	"395": "HeavySnowShowers",
}

// ===================================================================
// Emoji Profiles & Width Management (new centralized system)
// ===================================================================

// EmojiProfile defines a complete set of symbols + width strategy for alignment.
type EmojiProfile struct {
	Name    string
	Symbols map[string]string       // condition name → symbol
	Width   func(symbol string) int // returns terminal cells occupied (usually 1 or 2)
}

// Predefined profiles (used by the new ?emoji= option)
var EmojiProfiles = map[string]EmojiProfile{
	"unicode": { // default – standard Unicode emojis
		Name:    "unicode",
		Symbols: WeatherSymbol,
		// Width:   unicodeWidth, // most weather emojis are 2 cells
		Width: constantWidth(2),
	},
	"narrow": {
		Name:    "narrow",
		Symbols: WeatherSymbol,
		Width:   constantWidth(1),
	},
	"nerd": { // Weather Icons (PUA) – richest look, used by v2d/v2n
		Name:    "nerd",
		Symbols: WeatherSymbolWiDay, // can switch to WiNight based on time/view
		Width:   constantWidth(1),
	},
	"plain": { // ASCII fallback (used by StandardFont)
		Name:    "plain",
		Symbols: WeatherSymbolPlain,
		Width:   constantWidth(1),
	},
}

// Helper width functions
func constantWidth(w int) func(string) int {
	return func(string) int { return w }
}

func unicodeWidth(symbol string) int {
	if w, ok := WeatherSymbolWidthVTE[symbol]; ok {
		return w
	}
	return 2 // modern default for most weather emojis (Ghostty, Kitty, etc.)
}

// ===================================================================
// Legacy / Raw Maps (kept for backwards compatibility & specific use)
// ===================================================================

// WeatherSymbol maps symbolic condition names to Unicode emojis.
// Now built using full names via the enescakir/emoji library for maintainability.
var WeatherSymbol = map[string]string{
	"Unknown":             string(emoji.Sparkles),                  // ✨
	"Cloudy":              string(emoji.Cloud),                     // ☁️
	"Fog":                 string(emoji.Fog),                       // 🌫️
	"HeavyRain":           string(emoji.CloudWithRain),             // 🌧️
	"HeavyShowers":        string(emoji.CloudWithRain),             // 🌧️
	"HeavySnow":           string(emoji.CloudWithSnow),             // ❄️
	"HeavySnowShowers":    string(emoji.CloudWithSnow),             // ❄️
	"LightRain":           string(emoji.SunBehindRainCloud),        // 🌦️ (or CloudWithRain)
	"LightShowers":        string(emoji.SunBehindRainCloud),        // 🌦️
	"LightSleet":          string(emoji.CloudWithRain),             // 🌧️
	"LightSleetShowers":   string(emoji.CloudWithRain),             // 🌧️
	"LightSnow":           string(emoji.CloudWithSnow),             // 🌨️
	"LightSnowShowers":    string(emoji.CloudWithSnow),             // 🌨️
	"PartlyCloudy":        string(emoji.SunBehindSmallCloud),       // ⛅
	"Sunny":               string(emoji.Sun),                       // ☀️
	"ThunderyHeavyRain":   string(emoji.CloudWithLightning),        // 🌩️
	"ThunderyShowers":     string(emoji.CloudWithLightningAndRain), // ⛈️ (closest)
	"ThunderySnowShowers": string(emoji.CloudWithLightningAndRain), // ⛈️
	"VeryCloudy":          string(emoji.Cloud),                     // ☁️
}

// WeatherSymbolWidthVTE – legacy per-symbol width map (still used by unicode profile)
var WeatherSymbolWidthVTE = map[string]int{
	"✨": 2, "☁️": 1, "🌫️": 2, "🌧️": 2, "❄️": 1,
	"🌦️": 1, "🌨️": 2, "⛅": 2, "☀️": 1, "🌩️": 2, "⛈️": 1,
}

// WindDirection contains 8 cardinal arrow symbols (starting from North = ↓)
var WindDirection = [8]string{
	"↓", "↙", "←", "↖",
	"↑", "↗", "→", "↘",
}

// MoonPhases contains the 8 main moon phase emojis
var MoonPhases = [8]string{
	"🌑", "🌒", "🌓", "🌔",
	"🌕", "🌖", "🌗", "🌘",
}

// WeatherSymbolWiDay – Weather Icons style (day version)
// Unicode private use area characters (often rendered via custom font)
var WeatherSymbolWiDay = map[string]string{
	"Unknown":             "",
	"Cloudy":              "",
	"Fog":                 "",
	"HeavyRain":           "",
	"HeavyShowers":        "",
	"HeavySnow":           "",
	"HeavySnowShowers":    "",
	"LightRain":           "",
	"LightShowers":        "",
	"LightSleet":          "",
	"LightSleetShowers":   "",
	"LightSnow":           "",
	"LightSnowShowers":    "",
	"PartlyCloudy":        "",
	"Sunny":               "",
	"ThunderyHeavyRain":   "",
	"ThunderyShowers":     "",
	"ThunderySnowShowers": "",
	"VeryCloudy":          "",
}

// WeatherSymbolWiNight – Weather Icons style (night version)
var WeatherSymbolWiNight = map[string]string{
	"Unknown":             "",
	"Cloudy":              "",
	"Fog":                 "",
	"HeavyRain":           "",
	"HeavyShowers":        "",
	"HeavySnow":           "",
	"HeavySnowShowers":    "",
	"LightRain":           "",
	"LightShowers":        "",
	"LightSleet":          "",
	"LightSleetShowers":   "",
	"LightSnow":           "",
	"LightSnowShowers":    "",
	"PartlyCloudy":        "",
	"Sunny":               "",
	"ThunderyHeavyRain":   "",
	"ThunderyShowers":     "",
	"ThunderySnowShowers": "",
	"VeryCloudy":          "",
}

// WeatherSymbolPlain – very simple ASCII fallback representations
var WeatherSymbolPlain = map[string]string{
	"Unknown":             "?",
	"Cloudy":              "mm",
	"Fog":                 "=",
	"HeavyRain":           "///",
	"HeavyShowers":        "//",
	"HeavySnow":           "**",
	"HeavySnowShowers":    "*/*",
	"LightRain":           "/",
	"LightShowers":        ".",
	"LightSleet":          "x",
	"LightSleetShowers":   "x/",
	"LightSnow":           "*",
	"LightSnowShowers":    "*/",
	"PartlyCloudy":        "m",
	"Sunny":               "o",
	"ThunderyHeavyRain":   "/!/",
	"ThunderyShowers":     "!/",
	"ThunderySnowShowers": "*!*",
	"VeryCloudy":          "mmm",
}

// WindDirectionWi – Weather Icons style wind arrows
var WindDirectionWi = [8]string{
	"", "", "", "",
	"", "", "", "",
}

// WindScaleWi – Beaufort scale icons (0–12)
var WindScaleWi = [13]string{
	"", "", "", "", "",
	"", "", "", "", "",
	"", "", "",
}

// MoonPhasesWi – detailed moon phase icons (28 phases)
var MoonPhasesWi = [28]string{
	"", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "",
	"", "", "", "",
}

// LocaleToBCP47 maps short language codes to more complete locale identifiers
// (used when you need full locale strings for libraries or formatting)
var LocaleToBCP47 = map[string]string{
	"af":     "af_ZA",
	"ar":     "ar_TN",
	"az":     "az_AZ",
	"be":     "be_BY",
	"bg":     "bg_BG",
	"bs":     "bs_BA",
	"ca":     "ca_ES",
	"cs":     "cs_CZ",
	"cy":     "cy_GB",
	"da":     "da_DK",
	"de":     "de_DE",
	"el":     "el_GR",
	"eo":     "eo",
	"es":     "es_ES",
	"et":     "et_EE",
	"fa":     "fa_IR",
	"fi":     "fi_FI",
	"fr":     "fr_FR",
	"fy":     "fy_NL",
	"ga":     "ga_IE",
	"he":     "he_IL",
	"hr":     "hr_HR",
	"hu":     "hu_HU",
	"hy":     "hy_AM",
	"ia":     "ia",
	"id":     "id_ID",
	"is":     "is_IS",
	"it":     "it_IT",
	"ja":     "ja_JP",
	"jv":     "en_US", // fallback
	"ka":     "ka_GE",
	"ko":     "ko_KR",
	"kk":     "kk_KZ",
	"ky":     "ky_KG",
	"lt":     "lt_LT",
	"lv":     "lv_LV",
	"mk":     "mk_MK",
	"ml":     "ml_IN",
	"nb":     "nb_NO",
	"nl":     "nl_NL",
	"nn":     "nn_NO",
	"pt":     "pt_PT",
	"pt-br":  "pt_BR",
	"pl":     "pl_PL",
	"ro":     "ro_RO",
	"ru":     "ru_RU",
	"sv":     "sv_SE",
	"sk":     "sk_SK",
	"sl":     "sl_SI",
	"sr":     "sr_RS",
	"sr-lat": "sr_RS@latin",
	"sw":     "sw_KE",
	"th":     "th_TH",
	"tr":     "tr_TR",
	"uk":     "uk_UA",
	"uz":     "uz_UZ",
	"vi":     "vi_VN",
	"zh":     "zh_TW",
	"zu":     "zu_ZA",
	"mg":     "mg_MG",
}
