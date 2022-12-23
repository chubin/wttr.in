package v1

import (
	"math"
	"time"

	"github.com/klauspost/lctime"
)

var slotTimes = [slotcount]int{9 * 60, 12 * 60, 18 * 60, 22 * 60}

func printDay(w weather) (ret []string) {
	hourly := w.Hourly
	ret = make([]string, 5)
	for i := range ret {
		ret[i] = "│"
	}

	// find hourly data which fits the desired times of day best
	var slots [slotcount]cond
	for _, h := range hourly {
		c := int(math.Mod(float64(h.Time), 100)) + 60*(h.Time/100)
		for i, s := range slots {
			if math.Abs(float64(c-slotTimes[i])) < math.Abs(float64(s.Time-slotTimes[i])) {
				h.Time = c
				slots[i] = h
			}
		}
	}

	if config.RightToLeft {
		slots[0], slots[3] = slots[3], slots[0]
		slots[1], slots[2] = slots[2], slots[1]
	}

	for i, s := range slots {
		if config.Narrow {
			if i == 0 || i == 2 {
				continue
			}
		}
		ret = formatCond(ret, s, false)
		for i := range ret {
			ret[i] = ret[i] + "│"
		}
	}

	d, _ := time.Parse("2006-01-02", w.Date)
	// dateFmt := "┤ " + d.Format("Mon 02. Jan") + " ├"

	if val, ok := locale[config.Lang]; ok {
		lctime.SetLocale(val)
	} else {
		lctime.SetLocale("en_US")
	}
	dateName := ""
	if config.RightToLeft {
		dow := lctime.Strftime("%a", d)
		day := lctime.Strftime("%d", d)
		month := lctime.Strftime("%b", d)
		dateName = reverse(month) + " " + day + " " + reverse(dow)
	} else {
		dateName = lctime.Strftime("%a %d %b", d)
		if config.Lang == "ko" {
			dateName = lctime.Strftime("%b %d일 %a", d)
		}
		if config.Lang == "zh" || config.Lang == "zh-tw" || config.Lang == "zh-cn" {
			dateName = lctime.Strftime("%b%d日%A", d)
		}
	}
	// appendSide := 0
	// // for utf8.RuneCountInString(dateName) <= dateWidth {
	// for runewidth.StringWidth(dateName) <= dateWidth {
	//     if appendSide == 1 {
	//         dateName = dateName + " "
	//         appendSide = 0
	//     } else {
	//         dateName = " " + dateName
	//         appendSide = 1
	//     }
	// }

	dateFmt := "┤" + justifyCenter(dateName, 12) + "├"

	trans := daytimeTranslation["en"]
	if t, ok := daytimeTranslation[config.Lang]; ok {
		trans = t
	}
	if config.Narrow {

		names := "│      " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[3], 16) + "      │"

		ret = append([]string{
			"                        ┌─────────────┐                        ",
			"┌───────────────────────" + dateFmt + "───────────────────────┐",
			names,
			"├──────────────────────────────┼──────────────────────────────┤",
		},
			ret...)

		return append(ret,
			"└──────────────────────────────┴──────────────────────────────┘")

	}

	names := ""
	if config.RightToLeft {
		names = "│" + justifyCenter(trans[3], 29) + "│      " + justifyCenter(trans[2], 16) +
			"└──────┬──────┘" + justifyCenter(trans[1], 16) + "      │" + justifyCenter(trans[0], 29) + "│"
	} else {
		names = "│" + justifyCenter(trans[0], 29) + "│      " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[2], 16) + "      │" + justifyCenter(trans[3], 29) + "│"
	}

	ret = append([]string{
		"                                                       ┌─────────────┐                                                       ",
		"┌──────────────────────────────┬───────────────────────" + dateFmt + "───────────────────────┬──────────────────────────────┐",
		names,
		"├──────────────────────────────┼──────────────────────────────┼──────────────────────────────┼──────────────────────────────┤",
	},
		ret...)

	return append(ret,
		"└──────────────────────────────┴──────────────────────────────┴──────────────────────────────┴──────────────────────────────┘")
}
