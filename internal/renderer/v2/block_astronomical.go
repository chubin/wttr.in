package v2

import "github.com/chubin/wttr.in/internal/domain"

func drawAstronomical(loc *domain.Location) string {
	// Placeholder - can be expanded with real sunrise/sunset calculation later
	return "─ Sunrise ───── Noon ────── Sunset ───── Dusk ──\n" +
		"   06:12       13:05        20:58        22:10   \n"
}
