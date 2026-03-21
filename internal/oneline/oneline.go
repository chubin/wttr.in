// internal/oneline/oneline.go
package oneline

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chubin/wttr.go/internal/query"
	"github.com/chubin/wttr.go/internal/weather"
)

// ──────────────────────────────────────────────────────────────────────────────
// Unified render function signature
// ──────────────────────────────────────────────────────────────────────────────

type RenderFunc func(ctx *renderContext) string

// OnelineRenderer is responsible for rendering weather data
// in the compact one-line text format (?format=... or preconfigured IDs).
//
// It implements the weather.Renderer interface.
type OnelineRenderer struct {
	// Currently stateless — all configuration comes from query.Options
	// Future fields could include:
	// - symbolSets     map[string]map[int]string   // different emoji/icon sets
	// - moonCalculator MoonCalculator              // injectable moon phase logic
	// - astroProvider  AstroProvider               // sunrise/sunset calculation
}

// NewOnelineRenderer creates a new instance of the one-line renderer.
//
// At the moment it is stateless, so the constructor is very simple.
// It is kept as a constructor function to allow future dependency injection
// (e.g. custom symbol sets, astronomy calculators, caching, etc.).
func NewOnelineRenderer() *OnelineRenderer {
	return &OnelineRenderer{}
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers (usually kept private in the same file)
// ──────────────────────────────────────────────────────────────────────────────

func (r *OnelineRenderer) determineFormat(opts *query.Options) string {
	if opts == nil || opts.Format == "" {
		return "%c %t" // sensible minimal default
	}

	// Handle preconfigured format IDs from oneline.yaml
	switch opts.Format {
	case "1":
		return "%c %t\n"
	case "2":
		return "%c 🌡️%t 🌬️%w\n"
	case "3":
		return "%l: %c %t\n"
	case "4":
		return "%l: %c 🌡️%t 🌬️%w\n"
	case "69":
		return "nice"
	}

	// Otherwise treat as custom format string
	return opts.Format
}

// Example error definition (can be in the same package)
var ErrNoWeatherData = errors.New("no weather data available in query")

// renderContext holds everything a single placeholder might need
type renderContext struct {
	Data     *parsedCurrentCondition
	DataRaw  interface{}
	Options  *query.Options
	Location *weather.Location
	Now      time.Time
}

// ──────────────────────────────────────────────────────────────────────────────
// Central registry: letter → render function
// This matches format_specifiers.letter from oneline.yaml
// ──────────────────────────────────────────────────────────────────────────────

var placeholderRenderers = map[rune]RenderFunc{
	'c': renderConditionEmoji,
	'C': renderConditionFullName,
	'x': renderConditionPlain,
	'i': renderConditionCode,
	't': renderTemperature,
	'f': renderFeelsLike,
	'w': renderWind,
	'h': renderHumidity,
	'p': renderPrecipitation,
	'o': renderPrecipChance,
	'P': renderPressure,
	'u': renderUVIndex,
	'l': renderLocation,
	'm': renderMoonPhaseEmoji,
	'M': renderMoonDay,
	'S': renderSunrise,
	's': renderSunset,
	'D': renderDawn,
	'd': renderDusk,
	'z': renderSolarNoon,
	'T': renderLocalTime,
	'Z': renderTimezone,
}

// ──────────────────────────────────────────────────────────────────────────────
// Main rendering loop using the map
// ──────────────────────────────────────────────────────────────────────────────

func (r *OnelineRenderer) Render(q weather.Query) (weather.RenderOutput, error) {
	if q.Weather == nil || len(*q.Weather) == 0 {
		return weather.RenderOutput{}, fmt.Errorf("no weather data")
	}

	data, err := parseCurrentCondition(*q.Weather)
	if err != nil {
		return weather.RenderOutput{}, err
	}

	var dataRaw interface{} // or map[string]interface{}

	err = json.Unmarshal(*q.Weather, &dataRaw)
	if err != nil {
		return weather.RenderOutput{}, err
	}

	formatStr := r.determineFormat(q.Options)
	formatStr = TolerantUnescape(formatStr)

	ctx := &renderContext{
		Data:     data,
		DataRaw:  dataRaw,
		Options:  q.Options,
		Location: q.Location,
		Now:      time.Now(), // or from observation time if you prefer
	}

	var sb strings.Builder
	i := 0
	for i < len(formatStr) {
		if formatStr[i] == '%' && i+1 < len(formatStr) {
			letter := rune(formatStr[i+1])
			i += 2

			if fn, ok := placeholderRenderers[letter]; ok {
				sb.WriteString(fn(ctx))
			} else {
				// unknown placeholder → keep literal or log
				sb.WriteString(string([]rune{'%', letter}))
			}
			continue
		}
		sb.WriteByte(formatStr[i])
		i++
	}

	output := sb.String()

	// Replace \n in the text with new lines
	output = strings.ReplaceAll(output, `\n`, "\n")

	return weather.RenderOutput{
		Content: []byte(output),
	}, nil
}

func renderWithPlaceholders(format string, ctx *renderContext) string {
	var sb strings.Builder

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) {
			letter := format[i+1]
			i++ // skip the letter too

			if fn, ok := placeholderRenderers[rune(letter)]; ok {
				sb.WriteString(fn(ctx))
				continue
			}

			// Unknown → keep %X as-is
			sb.WriteByte('%')
			sb.WriteByte(letter)
			continue
		}

		sb.WriteByte(format[i])
	}

	return sb.String()
}

// TolerantUnescape decodes a string by replacing percent-encoded sequences (%XX)
// with their corresponding byte values, where XX is a two-digit hexadecimal number.
// If a percent-encoded sequence is invalid (e.g., not a valid hex number or incomplete),
// the original sequence is preserved as-is in the output. This makes the function
// "tolerant" to malformed input, unlike strict unescaping functions that might fail.
//
// For example:
//   - "%20" becomes a space character (ASCII 32).
//   - "%GG" remains "%GG" since "GG" is not a valid hex number.
//   - A lone "%" at the end of the string remains "%".
//
// Args:
//   s: The input string containing potential percent-encoded sequences.
//
// Returns:
//   A new string with valid percent-encoded sequences replaced by their decoded byte values.
func TolerantUnescape(s string) string {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '%' && i+2 < len(s) {
			hex := s[i+1 : i+3]
			if v, err := strconv.ParseUint(hex, 16, 8); err == nil {
				sb.WriteByte(byte(v))
				i += 3
				continue
			}
			// invalid hex → keep literal %XX
		}
		// lone % or invalid → keep as is
		sb.WriteByte(s[i])
		i++
	}
	return sb.String()
}
