//nolint:funlen,nestif,cyclop,gocognit
package v1

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/clipperhouse/displaywidth"

	"github.com/chubin/wttr.in/internal/options"
)

// windDir returns ANSI-colored arrow for wind direction.
func windDir() map[string]string {
	return map[string]string{
		"N":   "\033[1m↓\033[0m",
		"NNE": "\033[1m↓\033[0m",
		"NE":  "\033[1m↙\033[0m",
		"ENE": "\033[1m↙\033[0m",
		"E":   "\033[1m←\033[0m",
		"ESE": "\033[1m←\033[0m",
		"SE":  "\033[1m↖\033[0m",
		"SSE": "\033[1m↖\033[0m",
		"S":   "\033[1m↑\033[0m",
		"SSW": "\033[1m↑\033[0m",
		"SW":  "\033[1m↗\033[0m",
		"WSW": "\033[1m↗\033[0m",
		"W":   "\033[1m→\033[0m",
		"WNW": "\033[1m→\033[0m",
		"NW":  "\033[1m↘\033[0m",
		"NNW": "\033[1m↘\033[0m",
	}
}

// windInRightUnits converts wind speed to the requested unit (km/h, mph, or m/s).
func windInRightUnits(spd int, windMS, imperial bool) int {
	if windMS {
		return (spd * 1000) / 3600 // km/h → m/s
	}
	if imperial {
		return (spd * 1000) / 1609 // km/h → mph (approx)
	}
	return spd // default: km/h
}

// speedToColor returns colored wind speed text based on speed value.
func speedToColor(spd, spdConverted int) string {
	col := 46 // green
	switch {
	case spd <= 3:
		col = 82
	case spd <= 6:
		col = 118
	case spd <= 9:
		col = 154
	case spd <= 12:
		col = 190
	case spd <= 15:
		col = 226
	case spd <= 19:
		col = 220
	case spd <= 23:
		col = 214
	case spd <= 27:
		col = 208
	case spd <= 31:
		col = 202
	default:
		if spd > 0 {
			col = 196 // red
		}
	}
	return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, spdConverted)
}

// formatTemp formats temperature (with feels-like in parentheses when different).
func (r *V1Renderer) formatTemp(c cond, opts *options.Options) string {
	if opts == nil {
		opts = &options.Options{}
	}

	color := func(temp int, explicitPlus bool) string {
		var col int
		inverse := opts.InvertedColors

		if !inverse {
			col = 165
			switch {
			case temp >= -15 && temp <= -13:
				col = 171
			case temp >= -12 && temp <= -10:
				col = 33
			case temp >= -9 && temp <= -7:
				col = 39
			case temp >= -6 && temp <= -4:
				col = 45
			case temp >= -3 && temp <= -1:
				col = 51
			case temp == 0 || temp == 1:
				col = 50
			case temp == 2 || temp == 3:
				col = 49
			case temp == 4 || temp == 5:
				col = 48
			case temp == 6 || temp == 7:
				col = 47
			case temp == 8 || temp == 9:
				col = 46
			case temp >= 10 && temp <= 12:
				col = 82
			case temp >= 13 && temp <= 15:
				col = 118
			case temp >= 16 && temp <= 18:
				col = 154
			case temp >= 19 && temp <= 21:
				col = 190
			case temp >= 22 && temp <= 24:
				col = 226
			case temp >= 25 && temp <= 27:
				col = 220
			case temp >= 28 && temp <= 30:
				col = 214
			case temp >= 31 && temp <= 33:
				col = 208
			case temp >= 34 && temp <= 36:
				col = 202
			default:
				if temp > 0 {
					col = 196
				}
			}
		} else {
			col = 16
			switch {
			case temp >= -15 && temp <= -13:
				col = 17
			case temp >= -12 && temp <= -10:
				col = 18
			case temp >= -9 && temp <= -7:
				col = 19
			case temp >= -6 && temp <= -4:
				col = 20
			case temp >= -3 && temp <= -1:
				col = 21
			case temp == 0 || temp == 1:
				col = 30
			case temp == 2 || temp == 3:
				col = 28
			case temp == 4 || temp == 5:
				col = 29
			case temp == 6 || temp == 7:
				col = 30
			case temp == 8 || temp == 9:
				col = 34
			case temp >= 10 && temp <= 12:
				col = 35
			case temp >= 13 && temp <= 15:
				col = 36
			case temp >= 16 && temp <= 18:
				col = 40
			case temp >= 19 && temp <= 21:
				col = 59
			case temp >= 22 && temp <= 24:
				col = 100
			case temp >= 25 && temp <= 27:
				col = 101
			case temp >= 28 && temp <= 30:
				col = 94
			case temp >= 31 && temp <= 33:
				col = 166
			case temp >= 34 && temp <= 36:
				col = 52
			default:
				if temp > 0 {
					col = 196
				}
			}
		}

		if opts.UseImperial {
			temp = (temp*18 + 320) / 10
		}

		if explicitPlus {
			return fmt.Sprintf("\033[38;5;%03dm+%d\033[0m", col, temp)
		}
		return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, temp)
	}

	t := c.TempC
	if t == 0 {
		t = c.TempC2
	}

	explicitPlus1 := false
	explicitPlus2 := false
	if c.FeelsLikeC != t {
		if t > 0 {
			explicitPlus1 = true
		}
		if c.FeelsLikeC > 0 {
			explicitPlus2 = true
		}
		if explicitPlus1 {
			explicitPlus2 = false
		}

		return r.pad(
			fmt.Sprintf("%s(%s) °%s",
				color(t, explicitPlus1),
				color(c.FeelsLikeC, explicitPlus2),
				unitTemp()[opts.UseImperial]),
			15)
	}

	return r.pad(
		fmt.Sprintf("%s °%s", color(c.FeelsLikeC, false), unitTemp()[opts.UseImperial]),
		15)
}

// formatWind formats wind direction + speed (+ gust if stronger).
func (r *V1Renderer) formatWind(c cond, opts *options.Options) string {
	if opts == nil {
		opts = &options.Options{}
	}

	unitWindString := unitWind(0, opts.Lang)
	if opts.UseMsForWind {
		unitWindString = unitWind(2, opts.Lang)
	} else if opts.UseImperial {
		unitWindString = unitWind(1, opts.Lang)
	}

	gust := windInRightUnits(c.WindGustKmph, opts.UseMsForWind, opts.UseImperial)
	speed := windInRightUnits(c.WindspeedKmph, opts.UseMsForWind, opts.UseImperial)

	cWindGust := speedToColor(c.WindGustKmph, gust)
	cWindSpeed := speedToColor(c.WindspeedKmph, speed)

	hyphen := "-"

	if gust > speed {
		return r.pad(
			fmt.Sprintf("%s %s%s%s %s",
				windDir()[c.Winddir16Point],
				cWindSpeed, hyphen, cWindGust,
				unitWindString),
			15)
	}

	return r.pad(
		fmt.Sprintf("%s %s %s",
			windDir()[c.Winddir16Point],
			cWindSpeed,
			unitWindString),
		15)
}

// formatVisibility formats visibility.
func (r *V1Renderer) formatVisibility(c cond, opts *options.Options) string {
	if opts == nil {
		opts = &options.Options{}
	}

	vis := c.VisibleDistKM
	if opts.UseImperial {
		vis = (vis * 621) / 1000
	}

	return r.pad(
		fmt.Sprintf("%d %s", vis, unitVis(opts.UseImperial, opts.Lang)),
		15)
}

// formatRain formats precipitation amount and chance of rain.
func (r *V1Renderer) formatRain(c cond, opts *options.Options) string {
	if opts == nil {
		opts = &options.Options{}
	}

	rain := c.PrecipMM
	if opts.UseImperial {
		rain = c.PrecipMM * 0.039
	}

	if c.ChanceOfRain != "" {
		return r.pad(
			fmt.Sprintf("%.1f %s | %s%%",
				rain,
				unitRain(opts.UseImperial, opts.Lang),
				c.ChanceOfRain),
			15)
	}
	return r.pad(
		fmt.Sprintf("%.1f %s", rain, unitRain(opts.UseImperial, opts.Lang)),
		15)
}

// formatCond builds the 5-line current/forecast condition block.
// prefix is " " for current condition, "│" for forecast blocks inside the box.
func (r *V1Renderer) formatCond(prefix string, c cond, isCurrent bool, opts *options.Options) []string {
	if opts == nil {
		opts = &options.Options{}
	}
	if prefix == "" {
		prefix = " "
	}

	var ret []string

	// Get icon
	icon := getIcon("iconUnknown")
	if i, ok := codes()[c.WeatherCode]; ok {
		icon = i
	}

	// Inverse color adjustments
	if opts.InvertedColors {
		for i := range icon {
			icon[i] = strings.ReplaceAll(icon[i], "38;5;226", "38;5;94")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;250", "38;5;243")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;21", "38;5;18")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;255", "38;5;245")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;111", "38;5;63")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;251", "38;5;238")
		}
	}

	// Safe description handling
	desc := "Unknown"
	if len(c.WeatherDesc) > 0 {
		desc = c.WeatherDesc[0].Value
	}

	// Pad/truncate description to 15 characters
	if r.rightToLeft {
		for displaywidth.String(desc) < 15 {
			desc = " " + desc
		}
		for displaywidth.String(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[:len(desc)-size]
		}
	} else {
		for displaywidth.String(desc) < 15 {
			desc += " "
		}
		for displaywidth.String(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[:len(desc)-size]
		}
	}

	if isCurrent {
		if r.rightToLeft && displaywidth.String(desc) < 15 {
			desc = strings.Repeat(" ", 15-displaywidth.String(desc)) + desc
		} else {
			desc = strings.TrimRight(desc, " ")
		}
	} else {
		// Add ellipsis for forecast
		if r.rightToLeft {
			if first, size := utf8.DecodeRuneInString(desc); first != ' ' {
				desc = "…" + desc[size:]
				for displaywidth.String(desc) < 15 {
					desc = " " + desc
				}
			}
		} else {
			if last, size := utf8.DecodeLastRuneInString(desc); last != ' ' {
				desc = desc[:len(desc)-size] + "…"
				for displaywidth.String(desc) < 15 {
					desc += " "
				}
			}
		}
	}

	// Build the 5 lines using the prefix
	if r.rightToLeft {
		ret = append(ret,
			fmt.Sprintf("%s %s %s", prefix, desc, icon[0]),
			fmt.Sprintf("%s %s %s", prefix, r.formatTemp(c, opts), icon[1]),
			fmt.Sprintf("%s %s %s", prefix, r.formatWind(c, opts), icon[2]),
			fmt.Sprintf("%s %s %s", prefix, r.formatVisibility(c, opts), icon[3]),
			fmt.Sprintf("%s %s %s", prefix, r.formatRain(c, opts), icon[4]),
		)
	} else {
		ret = append(ret,
			fmt.Sprintf("%s %s %s", prefix, icon[0], desc),
			fmt.Sprintf("%s %s %s", prefix, icon[1], r.formatTemp(c, opts)),
			fmt.Sprintf("%s %s %s", prefix, icon[2], r.formatWind(c, opts)),
			fmt.Sprintf("%s %s %s", prefix, icon[3], r.formatVisibility(c, opts)),
			fmt.Sprintf("%s %s %s", prefix, icon[4], r.formatRain(c, opts)),
		)
	}

	return ret
}

// pad ensures the string has exactly `mustLen` visible runes, preserving ANSI codes.
func (r *V1Renderer) pad(s string, mustLen int) string {
	realLen := utf8.RuneCountInString(r.ansiEsc.ReplaceAllLiteralString(s, ""))
	delta := mustLen - realLen

	if delta > 0 {
		// Right-to-left support can be added later via opts if needed
		return s + "\033[0m" + strings.Repeat(" ", delta)
	}

	if delta < 0 {
		toks := r.ansiEsc.Split(s, 2)
		tokLen := utf8.RuneCountInString(toks[0])
		esc := r.ansiEsc.FindString(s)

		if tokLen > mustLen {
			return fmt.Sprintf("%.*s\033[0m", mustLen, toks[0])
		}
		return fmt.Sprintf("%s%s%s", toks[0], esc, r.pad(toks[1], mustLen-tokLen))
	}

	return s
}
