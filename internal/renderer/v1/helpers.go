package v1

import (
	"fmt"
	"strconv"

	"github.com/clipperhouse/displaywidth"
)

// Safe parsing helpers
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

// parseTimeToMinutes converts "0", "300", "600", ..., "2100" into minutes since midnight
func parseTimeToMinutes(timeInt int) int {
	timeStr := fmt.Sprintf("%04d", timeInt)
	var hour, minute int
	fmt.Sscanf(timeStr, "%2d%2d", &hour, &minute)
	return hour*60 + minute
}

// h2m returns minutes since midnight for a cond (used in comparison)
func h2m(h cond) int {
	return parseTimeToMinutes(h.Time)
}

func justifyCenter(s string, width int) string {
	appendSide := 0
	for displaywidth.String(s) <= width {
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
