package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

var (
	windDir = map[string]string{
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
)

func formatTemp(c cond) string {
	color := func(temp int, explicitPlus bool) string {
		var col = 0
		if !config.Inverse {
			// Extemely cold temperature must be shown with violet
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
		if config.Imperial {
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
		return pad(
			fmt.Sprintf("%s(%s) °%s",
				color(t, explicitPlus1),
				color(c.FeelsLikeC, explicitPlus2),
				unitTemp[config.Imperial]),
			15)
	}
	// if c.FeelsLikeC < t {
	// 	if c.FeelsLikeC < 0 && t > 0 {
	// 		explicitPlus = true
	// 	}
	// 	return pad(fmt.Sprintf("%s%s%s °%s", color(c.FeelsLikeC, false), hyphen, color(t, explicitPlus), unitTemp[config.Imperial]), 15)
	// } else if c.FeelsLikeC > t {
	// 	if t < 0 && c.FeelsLikeC > 0 {
	// 		explicitPlus = true
	// 	}
	// 	return pad(fmt.Sprintf("%s%s%s °%s", color(t, false), hyphen, color(c.FeelsLikeC, explicitPlus), unitTemp[config.Imperial]), 15)
	// }
	return pad(fmt.Sprintf("%s °%s", color(c.FeelsLikeC, false), unitTemp[config.Imperial]), 15)
}

func formatWind(c cond) string {
	windInRightUnits := func(spd int) int {
		if config.WindMS {
			spd = (spd * 1000) / 3600
		} else {
			if config.Imperial {
				spd = (spd * 1000) / 1609
			}
		}
		return spd
	}
	color := func(spd int) string {
		var col = 46
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
		spd = windInRightUnits(spd)

		return fmt.Sprintf("\033[38;5;%03dm%d\033[0m", col, spd)
	}

	unitWindString := unitWind(0, config.Lang)
	if config.WindMS {
		unitWindString = unitWind(2, config.Lang)
	} else {
		if config.Imperial {
			unitWindString = unitWind(1, config.Lang)
		}
	}

	hyphen := " - "
	// if (config.Lang == "sl") {
	//     hyphen = "-"
	// }
	hyphen = "-"

	cWindGustKmph := color(c.WindGustKmph)
	cWindspeedKmph := color(c.WindspeedKmph)
	if windInRightUnits(c.WindGustKmph) > windInRightUnits(c.WindspeedKmph) {
		return pad(fmt.Sprintf("%s %s%s%s %s", windDir[c.Winddir16Point], cWindspeedKmph, hyphen, cWindGustKmph, unitWindString), 15)
	}
	return pad(fmt.Sprintf("%s %s %s", windDir[c.Winddir16Point], cWindspeedKmph, unitWindString), 15)
}

func formatVisibility(c cond) string {
	if config.Imperial {
		c.VisibleDistKM = (c.VisibleDistKM * 621) / 1000
	}
	return pad(fmt.Sprintf("%d %s", c.VisibleDistKM, unitVis(config.Imperial, config.Lang)), 15)
}

func formatRain(c cond) string {
	rainUnit := float32(c.PrecipMM)
	if config.Imperial {
		rainUnit = float32(c.PrecipMM) * 0.039
	}
	if c.ChanceOfRain != "" {
		return pad(fmt.Sprintf(
			"%.1f %s | %s%%",
			rainUnit,
			unitRain(config.Imperial, config.Lang),
			c.ChanceOfRain), 15)
	}
	return pad(fmt.Sprintf("%.1f %s", rainUnit, unitRain(config.Imperial, config.Lang)), 15)
}

func formatCond(cur []string, c cond, current bool) (ret []string) {
	var icon []string
	if i, ok := codes[c.WeatherCode]; !ok {
		icon = iconUnknown
	} else {
		icon = i
	}
	if config.Inverse {
		// inverting colors
		for i := range icon {
			icon[i] = strings.Replace(icon[i], "38;5;226", "38;5;94", -1)
			icon[i] = strings.Replace(icon[i], "38;5;250", "38;5;243", -1)
			icon[i] = strings.Replace(icon[i], "38;5;21", "38;5;18", -1)
			icon[i] = strings.Replace(icon[i], "38;5;255", "38;5;245", -1)
			icon[i] = strings.Replace(icon[i], "38;5;111", "38;5;63", -1)
			icon[i] = strings.Replace(icon[i], "38;5;251", "38;5;238", -1)
		}
	}
	//desc := fmt.Sprintf("%-15.15v", c.WeatherDesc[0].Value)
	desc := c.WeatherDesc[0].Value
	if config.RightToLeft {
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
		if config.RightToLeft {
			desc = c.WeatherDesc[0].Value
			if runewidth.StringWidth(desc) < 15 {
				desc = strings.Repeat(" ", 15-runewidth.StringWidth(desc)) + desc
			}
		} else {
			desc = c.WeatherDesc[0].Value
		}
	} else {
		if config.RightToLeft {
			if frstRune, size := utf8.DecodeRuneInString(desc); frstRune != ' ' {
				desc = "…" + desc[size:]
				for runewidth.StringWidth(desc) < 15 {
					desc = " " + desc
				}
			}
		} else {
			if lastRune, size := utf8.DecodeLastRuneInString(desc); lastRune != ' ' {
				desc = desc[:len(desc)-size] + "…"
				//for numberOfSpaces < runewidth.StringWidth(fmt.Sprintf("%c", lastRune)) - 1 {
				for runewidth.StringWidth(desc) < 15 {
					desc = desc + " "
				}
			}
		}
	}
	if config.RightToLeft {
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[0], desc, icon[0]), fmt.Sprintf("%v %v %v", cur[1], formatTemp(c), icon[1]), fmt.Sprintf("%v %v %v", cur[2], formatWind(c), icon[2]), fmt.Sprintf("%v %v %v", cur[3], formatVisibility(c), icon[3]), fmt.Sprintf("%v %v %v", cur[4], formatRain(c), icon[4]))
	} else {
		ret = append(ret, fmt.Sprintf("%v %v %v", cur[0], icon[0], desc), fmt.Sprintf("%v %v %v", cur[1], icon[1], formatTemp(c)), fmt.Sprintf("%v %v %v", cur[2], icon[2], formatWind(c)), fmt.Sprintf("%v %v %v", cur[3], icon[3], formatVisibility(c)), fmt.Sprintf("%v %v %v", cur[4], icon[4], formatRain(c)))
	}
	return
}

func justifyCenter(s string, width int) string {
	appendSide := 0
	for runewidth.StringWidth(s) <= width {
		if appendSide == 1 {
			s = s + " "
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

func pad(s string, mustLen int) (ret string) {
	ret = s
	realLen := utf8.RuneCountInString(ansiEsc.ReplaceAllLiteralString(s, ""))
	delta := mustLen - realLen
	if delta > 0 {
		if config.RightToLeft {
			ret = strings.Repeat(" ", delta) + ret + "\033[0m"
		} else {
			ret += "\033[0m" + strings.Repeat(" ", delta)
		}
	} else if delta < 0 {
		toks := ansiEsc.Split(s, 2)
		tokLen := utf8.RuneCountInString(toks[0])
		esc := ansiEsc.FindString(s)
		if tokLen > mustLen {
			ret = fmt.Sprintf("%.*s\033[0m", mustLen, toks[0])
		} else {
			ret = fmt.Sprintf("%s%s%s", toks[0], esc, pad(toks[1], mustLen-tokLen))
		}
	}
	return
}
