package oneline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// RenderConditionEmoji returns emoji/symbol + appropriate spacing
// Behaves as closely as possible to the original Python version:
//
//   if view == "v2n":
//       symbol = WeatherSymbolWiNight[WWO_CODE.get(code, "Unknown")]
//       spaces = " "
//   elif view == "v2d":
//       symbol = WeatherSymbolWiDay[WWO_CODE.get(code, "Unknown")]
//       spaces = " "
//   else:
//       symbol = WeatherSymbol.get(WWO_CODE.get(code, "Unknown"))
//       spaces = " " * (3 - WeatherSymbolWidthVTE.get(symbol, 1))
//
// If code is missing/invalid → warning to stderr + placeholder symbol
func RenderConditionEmoji(ctx *RenderContext) string {
	if ctx == nil || ctx.Data == nil || ctx.Data.ConditionCode == "" {
		fmt.Fprintln(os.Stderr, "WARNING: renderCondition called with empty or nil ConditionCode")
		return "❓"
	}

	codeStr := ctx.Data.ConditionCode

	// Try to convert to int just to validate it's numeric
	_, err := strconv.Atoi(codeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: invalid weather code (not numeric): %q\n", codeStr)
		return "❓"
	}

	view := ""
	if ctx.Options != nil {
		view = ctx.Options.View
	}

	// ────────────────────────────────────────────────────────────────
	// Get internal symbolic name (same as Python's WWO_CODE.get(...))
	// ────────────────────────────────────────────────────────────────
	name, ok := WWOCodeToName[codeStr]
	if !ok {
		fmt.Fprintf(os.Stderr, "WARNING: unknown weather code %q – not in WWOCodeToName\n", codeStr)
		name = "Unknown"
	}

	// ────────────────────────────────────────────────────────────────
	// v2n → night Weather Icons
	// ────────────────────────────────────────────────────────────────
	if view == "v2n" {
		symbol, found := WeatherSymbolWiNight[name]
		if !found {
			fmt.Fprintf(os.Stderr, "WARNING: no night Wi symbol for %q\n", name)
			return "" // same as Wi Unknown
		}
		return symbol + " "
	}

	// ────────────────────────────────────────────────────────────────
	// v2d → day Weather Icons
	// ────────────────────────────────────────────────────────────────
	if view == "v2d" {
		symbol, found := WeatherSymbolWiDay[name]
		if !found {
			fmt.Fprintf(os.Stderr, "WARNING: no day Wi symbol for %q\n", name)
			return ""
		}
		return symbol + " "
	}

	// ────────────────────────────────────────────────────────────────
	// default → colorful emoji + calculated padding
	// ────────────────────────────────────────────────────────────────
	symbol, found := WeatherSymbol[name]
	if !found {
		fmt.Fprintf(os.Stderr, "WARNING: no emoji symbol defined for %q (code %s)\n", name, codeStr)
		return "✨"
	}

	// Padding logic exactly like Python
	width := WeatherSymbolWidthVTE[symbol]
	if width == 0 {
		width = 1 // default like in Python .get(..., 1)
	}
	padding := 3 - width
	if padding < 0 {
		padding = 0
	}

	return symbol + strings.Repeat(" ", padding)
}
