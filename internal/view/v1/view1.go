package v1

import (
	"math"
	"time"

	"github.com/klauspost/lctime"
)

func slotTimes() []int {
	return []int{9 * 60, 12 * 60, 18 * 60, 22 * 60}
}

//nolint:funlen,gocognit,cyclop
func (g *global) printDay(w weather) ([]string, error) {
	var (
		ret      = []string{}
		dateName string
		names    string
	)

	hourly := w.Hourly
	for i := range ret {
		ret[i] = "│"
	}

	// find hourly data which fits the desired times of day best
	var slots [slotcount]cond
	for _, h := range hourly {
		c := int(math.Mod(float64(h.Time), 100)) + 60*(h.Time/100)
		for i, s := range slots {
			if math.Abs(float64(c-slotTimes()[i])) < math.Abs(float64(s.Time-slotTimes()[i])) {
				h.Time = c
				slots[i] = h
			}
		}
	}

	if g.config.RightToLeft {
		slots[0], slots[3] = slots[3], slots[0]
		slots[1], slots[2] = slots[2], slots[1]
	}

	for i, s := range slots {
		if g.config.Narrow {
			if i == 0 || i == 2 {
				continue
			}
		}
		ret = g.formatCond(ret, s, false)
		for i := range ret {
			ret[i] += "│"
		}
	}

	d, _ := time.Parse("2006-01-02", w.Date)
	// dateFmt := "┤ " + d.Format("Mon 02. Jan") + " ├"

	if val, ok := locale()[g.config.Lang]; ok {
		err := lctime.SetLocale(val)
		if err != nil {
			return nil, err
		}
	} else {
		err := lctime.SetLocale("en_US")
		if err != nil {
			return nil, err
		}
	}

	if g.config.RightToLeft {
		dow := lctime.Strftime("%a", d)
		day := lctime.Strftime("%d", d)
		month := lctime.Strftime("%b", d)
		dateName = reverse(month) + " " + day + " " + reverse(dow)
	} else {
		switch g.config.Lang {
		case "ko":
			date_format = "%b %d일 %a"
		case "lv":
			date_format = "%a., %d. %b."
		case "zh", "zh-cn", "zh-tw":
			date_format = "%b%d日%A"
		default:
			date_format = "%a %d %b"
		}
		dateName = lctime.Strftime(date_format, d)
	}

	dateFmt := "┤" + justifyCenter(dateName, 12) + "├"

	trans := daytimeTranslation()["en"]
	if t, ok := daytimeTranslation()[g.config.Lang]; ok {
		trans = t
	}
	if g.config.Narrow {
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
				"└──────────────────────────────┴──────────────────────────────┘"),
			nil
	}

	if g.config.RightToLeft {
		names = "│" + justifyCenter(trans[3], 29) + "│      " + justifyCenter(trans[2], 16) +
			"└──────┬──────┘" + justifyCenter(trans[1], 16) + "      │" + justifyCenter(trans[0], 29) + "│"
	} else {
		names = "│" + justifyCenter(trans[0], 29) + "│      " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[2], 16) + "      │" + justifyCenter(trans[3], 29) + "│"
	}

	//nolint:lll
	ret = append([]string{
		"                                                       ┌─────────────┐                                                       ",
		"┌──────────────────────────────┬───────────────────────" + dateFmt + "───────────────────────┬──────────────────────────────┐",
		names,
		"├──────────────────────────────┼──────────────────────────────┼──────────────────────────────┼──────────────────────────────┤",
	},
		ret...)

	//nolint:lll
	return append(ret,
			"└──────────────────────────────┴──────────────────────────────┴──────────────────────────────┴──────────────────────────────┘"),
		nil
}
