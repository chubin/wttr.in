//nolint:funlen,nestif,cyclop,gocognit,gocyclo
package v1

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

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

func (g *global) formatTemp(c cond) string {
	color := func(temp int, explicitPlus bool) string {
		var col int
		//nolint:dupl
		if !g.config.Inverse {
			// Extremely cold temperature must be shown with violet
			// because dark blue is too dark
			col = 165
			switch temp {
			case -15, -14, -13:
				col = 171
			case -12, -11, -10:
				col = 33
			case -9, -8, -7:
				col = 39
			case -6, -5, -4:
				col = 45
			case -3, -2, -1:
				col = 51
			case 0, 1:
				col = 50
			case 2, 3:
				col = 49
			case 4, 5:
				col = 48
			case 6, 7:
				col = 47
			case 8, 9:
				col = 46
			case 10, 11, 12:
				col = 82
			case 13, 14, 15:
				col = 118
			case 16, 17, 18:
				col = 154
			case 19, 20, 21:
				col = 190
			case 22, 23, 24:
				col = 226
			case 25, 26, 27:
				col = 220
			case 28, 29, 30:
				col = 214
			case 31, 32, 33:
				col = 208
			case 34, 35, 36:
				col = 202
			default:
				if temp > 0 {
					col = 196
				}
			}
		} else {
			col = 16
			switch temp {
			case -15, -14, -13:
				col = 17
			case -12, -11, -10:
				col = 18
			case -9, -8, -7:
				col = 19
			case -6, -5, -4:
				col = 20
			case -3, -2, -1:
				col = 21
			case 0, 1:
				col = 30
			case 2, 3:
				col = 28
			case 4, 5:
				col = 29
			case 6, 7:
				col = 30
			case 8, 9:
				col = 34
			case 10, 11, 12:
				col = 35
			case 13, 14, 15:
				col = 36
			case 16, 17, 18:
				col = 40
			case 19, 20, 21:
				col = 59
			case 22, 23, 24:
				col = 100
			case 25, 26, 27:
				col = 101
			case 28, 29, 30:
				col = 94
			case 31, 32, 33:
				col = 166
			case 34, 35, 36:
				col = 52
			default:
				if temp > 0 {
					col = 196
				}
			}
		}
		if g.config.Imperial {
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

	// hyphen := " - "

	// if (config.Lang == "sl") {
	//     hyphen = "-"
	// }

	// hyphen = ".."

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

		return g.pad(
			fmt.Sprintf("%s(%s) °%s",
				color(t, explicitPlus1),
				color(c.FeelsLikeC, explicitPlus2),
				unitTemp()[g.config.Imperial]),
			15)
	}

	return g.pad(fmt.Sprintf("%s °%s", color(c.FeelsLikeC, false), unitTemp()[g.config.Imperial]), 15)
}

func (g *global) formatWind(c cond) string {
	unitWindString := unitWind(0, g.config.Lang)
	if g.config.WindMS {
		unitWindString = unitWind(2, g.config.Lang)
	} else if g.config.Imperial {
		unitWindString = unitWind(1, g.config.Lang)
	}

	hyphen := "-"

	cWindGustKmph := speedToColor(c.WindGustKmph, windInRightUnits(c.WindGustKmph, g.config.WindMS, g.config.Imperial))
	cWindspeedKmph := speedToColor(c.WindspeedKmph, windInRightUnits(c.WindspeedKmph, g.config.WindMS, g.config.Imperial))
	if windInRightUnits(c.WindGustKmph, g.config.WindMS, g.config.Imperial) >
		windInRightUnits(c.WindspeedKmph, g.config.WindMS, g.config.Imperial) {
		return g.pad(
			fmt.Sprintf("%s %s%s%s %s", windDir()[c.Winddir16Point], cWindspeedKmph, hyphen, cWindGustKmph, unitWindString),
			15)
	}

	return g.pad(fmt.Sprintf("%s %s %s", windDir()[c.Winddir16Point], cWindspeedKmph, unitWindString), 15)
}

func windInRightUnits(spd int, windMS, imperial bool) int {
	if windMS {
		spd = (spd * 1000) / 3600
	} else if imperial {
		spd = (spd * 1000) / 1609
	}

	return spd
}

func speedToColor(spd, spdConverted int) string {
	col := 46
	switch spd {
	case 1, 2, 3:
		col = 82
	case 4, 5, 6:
		col = 118
	case 7, 8, 9:
		col = 154
	case 10, 11, 12:
		col = 190
	case 13, 14, 15:
		col = 226
	case 16, 17, 18, 19:
		col = 220
	case 20, 21, 22, 23:
		col = 214
	case 24, 25, 26, 27:
		col = 208
	case 28, 29, 30, 31:
		col = 202
	default:
		if spd > 0 {
			col = 196
		}
	}

	return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, spdConverted)
}

func (g *global) formatVisibility(c cond) string {
	if g.config.Imperial {
		c.VisibleDistKM = (c.VisibleDistKM * 621) / 1000
	}

	return g.pad(fmt.Sprintf("%d %s", c.VisibleDistKM, unitVis(g.config.Imperial, g.config.Lang)), 15)
}

func (g *global) formatRain(c cond) string {
	rainUnit := c.PrecipMM
	if g.config.Imperial {
		rainUnit = c.PrecipMM * 0.039
	}
	if c.ChanceOfRain != "" {
		return g.pad(fmt.Sprintf(
			"%.1f %s | %s%%",
			rainUnit,
			unitRain(g.config.Imperial, g.config.Lang),
			c.ChanceOfRain), 15)
	}

	return g.pad(fmt.Sprintf("%.1f %s", rainUnit, unitRain(g.config.Imperial, g.config.Lang)), 15)
}

func (g *global) formatCond(cur []string, c cond, current bool) []string {
	var (
		ret  []string
		icon []string
	)

	if i, ok := codes()[c.WeatherCode]; !ok {
		icon = getIcon("iconUnknown")
	} else {
		icon = i
	}
	if g.config.Inverse {
		// inverting colors
		for i := range icon {
			icon[i] = strings.ReplaceAll(icon[i], "38;5;226", "38;5;94")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;250", "38;5;243")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;21", "38;5;18")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;255", "38;5;245")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;111", "38;5;63")
			icon[i] = strings.ReplaceAll(icon[i], "38;5;251", "38;5;238")
		}
	}
	// desc := fmt.Sprintf("%-15.15v", c.WeatherDesc[0].Value)
	desc := c.WeatherDesc[0].Value
	if g.config.RightToLeft {
		for runewidth.StringWidth(desc) < 15 {
			desc = " " + desc
		}
		for runewidth.StringWidth(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[size:]
		}
	} else {
		for runewidth.StringWidth(desc) < 15 {
			desc += " "
		}
		for runewidth.StringWidth(desc) > 15 {
			_, size := utf8.DecodeLastRuneInString(desc)
			desc = desc[:len(desc)-size]
		}
	}
	if current {
		if g.config.RightToLeft {
			desc = c.WeatherDesc[0].Value
			if runewidth.StringWidth(desc) < 15 {
				desc = strings.Repeat(" ", 15-runewidth.StringWidth(desc)) + desc
			}
		} else {
			desc = c.WeatherDesc[0].Value
		}
	} else {
		if g.config.RightToLeft {
			if frstRune, size := utf8.DecodeRuneInString(desc); frstRune != ' ' {
				desc = "…" + desc[size:]
				for runewidth.StringWidth(desc) < 15 {
					desc = " " + desc
				}
			}
		} else {
			if lastRune, size := utf8.DecodeLastRuneInString(desc); lastRune != ' ' {
				desc = desc[:len(desc)-size] + "…"
				// for numberOfSpaces < runewidth.StringWidth(fmt.Sprintf("%c", lastRune)) - 1 {
				for runewidth.StringWidth(desc) < 15 {
					desc += " "
				}
			}
		}
	}
	if g.config.RightToLeft {
		ret = append(
			ret,
			fmt.Sprintf("%v %v %v", cur[0], desc, icon[0]),
			fmt.Sprintf("%v %v %v", cur[1], g.formatTemp(c), icon[1]),
			fmt.Sprintf("%v %v %v", cur[2], g.formatWind(c), icon[2]),
			fmt.Sprintf("%v %v %v", cur[3], g.formatVisibility(c), icon[3]),
			fmt.Sprintf("%v %v %v", cur[4], g.formatRain(c), icon[4]))
	} else {
		ret = append(
			ret,
			fmt.Sprintf("%v %v %v", cur[0], icon[0], desc),
			fmt.Sprintf("%v %v %v", cur[1], icon[1], g.formatTemp(c)),
			fmt.Sprintf("%v %v %v", cur[2], icon[2], g.formatWind(c)),
			fmt.Sprintf("%v %v %v", cur[3], icon[3], g.formatVisibility(c)),
			fmt.Sprintf("%v %v %v", cur[4], icon[4], g.formatRain(c)))
	}

	return ret
}

func justifyCenter(s string, width int) string {
	appendSide := 0
	for runewidth.StringWidth(s) <= width {
		if appendSide == 1 {
			s += " "
			appendSide = 0
		} else {
			s = " " + s
			appendSide = 1
		}
	}

	return s
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

func (g *global) pad(s string, mustLen int) string {
	var ret string
	ret = s
	realLen := utf8.RuneCountInString(g.ansiEsc.ReplaceAllLiteralString(s, ""))
	delta := mustLen - realLen
	if delta > 0 {
		if g.config.RightToLeft {
			ret = strings.Repeat(" ", delta) + ret + "\033[0m"
		} else {
			ret += "\033[0m" + strings.Repeat(" ", delta)
		}
	} else if delta < 0 {
		toks := g.ansiEsc.Split(s, 2)
		tokLen := utf8.RuneCountInString(toks[0])
		esc := g.ansiEsc.FindString(s)
		if tokLen > mustLen {
			ret = fmt.Sprintf("%.*s\033[0m", mustLen, toks[0])
		} else {
			ret = fmt.Sprintf("%s%s%s", toks[0], esc, g.pad(toks[1], mustLen-tokLen))
		}
	}

	return ret
}
