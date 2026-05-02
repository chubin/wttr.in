package v2

import (
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
)

// drawDate - respects NoCaption / NoCity indirectly through frame
func drawDate(loc *domain.Location) string {
	tz, _ := time.LoadLocation(loc.TimeZone)

	// Match Python: start from UTC and add full days (more consistent across timezones)
	now := time.Now().UTC()
	tzNow := now.In(tz)

	var b strings.Builder

	// 3 date lines
	for i := 0; i < 3; i++ {
		d := tzNow.AddDate(0, 0, i)
		dateStr := d.Format("Mon 02 Jan")

		// Center in 24 characters (Python-style padding: extra space on left if odd)
		total := 24
		left := (total - len(dateStr)) / 2
		right := total - len(dateStr) - left

		// Python does left-heavy when odd length, but difference is tiny
		line := strings.Repeat(" ", left) + dateStr + strings.Repeat(" ", right)
		b.WriteString(line)
	}
	b.WriteString("\n")

	// 3 separator lines with ╷
	for i := 0; i < 3; i++ {
		tick := "╷"
		if i == 2 {
			tick = " "
		}
		b.WriteString(strings.Repeat(" ", 23) + tick)
	}
	b.WriteString("\n")

	result := b.String()

	return result
}
