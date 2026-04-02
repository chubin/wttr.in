package v1

import (
	"fmt"
	"math"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/klauspost/lctime"
	"github.com/mattn/go-runewidth"
)

// slotTimes returns the preferred times of day (in minutes since midnight) for the 4 slots.
func slotTimes() []int {
	return []int{9 * 60, 12 * 60, 18 * 60, 22 * 60}
}

// printDay renders a single forecast day in the classic v1 box style.
func (r *V1Renderer) printDay(day domain.WeatherDay, opts *options.Options) ([]string, error) {
	if opts == nil {
		opts = &options.Options{}
	}

	// Select best hourly slots for morning, noon, evening, night
	slots := r.selectBestHourlySlots(day.Hourly)

	// Handle right-to-left languages
	if r.rightToLeft {
		slots[0], slots[3] = slots[3], slots[0]
		slots[1], slots[2] = slots[2], slots[1]
	}

	// Format the four hourly blocks
	var ret []string
	for _, slot := range slots {
		if opts.Narrow && (len(ret) == 1 || len(ret) == 3) { // skip some slots in narrow mode
			continue
		}
		lines := r.formatHourlyCondition(slot, false, opts)
		ret = append(ret, lines...)
	}

	// Build date header with localization
	dateName, err := r.formatDate(day.Date, opts)
	if err != nil {
		return nil, err
	}
	dateFmt := "┤" + justifyCenter(dateName, 12) + "├"

	// Daytime translations (Morning, Noon, Evening, Night)
	trans := daytimeTranslation()["en"]
	if t, ok := daytimeTranslation()[opts.Lang]; ok {
		trans = t
	}

	// Narrow mode layout
	if opts.Narrow {
		names := "│ " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[3], 16) + " │"

		ret = append([]string{
			" ┌─────────────┐ ",
			"┌───────────────────────" + dateFmt + "───────────────────────┐",
			names,
			"├──────────────────────────────┼──────────────────────────────┤",
		}, ret...)

		ret = append(ret,
			"└──────────────────────────────┴──────────────────────────────┘",
		)
		return ret, nil
	}

	// Normal (wide) mode layout
	var names string
	if r.rightToLeft {
		names = "│" + justifyCenter(trans[3], 29) + "│ " + justifyCenter(trans[2], 16) +
			"└──────┬──────┘" + justifyCenter(trans[1], 16) + " │" + justifyCenter(trans[0], 29) + "│"
	} else {
		names = "│" + justifyCenter(trans[0], 29) + "│ " + justifyCenter(trans[1], 16) +
			"└──────┬──────┘" + justifyCenter(trans[2], 16) + " │" + justifyCenter(trans[3], 29) + "│"
	}

	ret = append([]string{
		" ┌─────────────┐ ",
		"┌──────────────────────────────┬───────────────────────" + dateFmt + "───────────────────────┬──────────────────────────────┐",
		names,
		"├──────────────────────────────┼──────────────────────────────┼──────────────────────────────┼──────────────────────────────┤",
	}, ret...)

	ret = append(ret,
		"└──────────────────────────────┴──────────────────────────────┴──────────────────────────────┴──────────────────────────────┘",
	)

	return ret, nil
}

// selectBestHourlySlots picks the best matching hourly entry for each target time slot.
func (r *V1Renderer) selectBestHourlySlots(hourly []domain.Hourly) [4]domain.Hourly {
	var slots [4]domain.Hourly
	targets := slotTimes()

	for _, h := range hourly {
		// Convert "HHMM" string like "0900", "1200" to minutes since midnight
		hour := 0
		minute := 0
		fmt.Sscanf(h.Time, "%2d%2d", &hour, &minute)
		minutes := hour*60 + minute

		for i, target := range targets {
			if math.Abs(float64(minutes-target)) < math.Abs(float64(h2m(slots[i])-target)) {
				slots[i] = h
			}
		}
	}
	return slots
}

// h2m returns the time of an Hourly entry in minutes since midnight.
// Helper to make the comparison cleaner.
func h2m(h domain.Hourly) int {
	if h.Time == "" {
		return 0
	}
	hour := 0
	minute := 0
	fmt.Sscanf(h.Time, "%2d%2d", &hour, &minute)
	return hour*60 + minute
}

// formatHourlyCondition renders one hourly block (morning/noon/evening/night).
// The `isCurrent` flag is not used here (kept for compatibility with formatCond).
func (r *V1Renderer) formatHourlyCondition(h domain.Hourly, isCurrent bool, opts *options.Options) []string {
	// TODO: Implement full classic formatting with icon, temperature, wind, etc.
	// For now we provide a minimal working version that matches the structure.

	temp := h.TempC
	if opts.UseImperial {
		temp = h.TempF
	}

	icon := getWeatherIcon(h.WeatherCode) // you should implement this

	line1 := fmt.Sprintf(" %s %s°", icon, temp)
	line2 := fmt.Sprintf(" %s", h.WeatherDesc[0].Value)

	return []string{line1, line2}
}

// formatDate returns localized and optionally reversed date string for the header.
func (r *V1Renderer) formatDate(dateStr string, opts *options.Options) (string, error) {
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date %s: %w", dateStr, err)
	}

	// Set locale for lctime
	localeStr := "en_US"
	if val, ok := locale()[opts.Lang]; ok {
		localeStr = val
	}
	if err := lctime.SetLocale(localeStr); err != nil {
		return "", err
	}

	var dateName string
	if r.rightToLeft {
		dow := lctime.Strftime("%a", d)
		day := lctime.Strftime("%d", d)
		month := lctime.Strftime("%b", d)
		dateName = reverse(month) + " " + day + " " + reverse(dow)
	} else {
		dateFormat := "%a %d %b"
		switch opts.Lang {
		case "ko":
			dateFormat = "%b %d일 %a"
		case "lv":
			dateFormat = "%a., %d. %b."
		case "zh", "zh-cn", "zh-tw":
			dateFormat = "%b%d日%A"
		}
		dateName = lctime.Strftime(dateFormat, d)
	}

	return dateName, nil
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

func getWeatherIcon(c string) string {
	return ""
}
