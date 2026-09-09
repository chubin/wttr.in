package v2

import (
	"strings"
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
)

// TestBuildLocationString_RTLMark checks that the location line gets a
// leading right-to-left mark (U+200F) for RTL languages, and stays
// plain otherwise. Without it, terminals display RTL location names
// in the wrong visual order (issue #932).
func TestBuildLocationString_RTLMark(t *testing.T) {
	loc := &domain.Location{
		Name:        "Talagang",
		Country:     "Pakistan",
		FullAddress: "تلہ گنگ, ضلع اٹک, پنجاب, پاکستان",
		Latitude:    32.9235108,
		Longitude:   72.4139769,
	}

	ltr := buildLocationString(loc, false)
	if strings.HasPrefix(ltr, rlm) {
		t.Errorf("buildLocationString(isRTL=false) should not start with RLM, got %q", ltr)
	}
	if !strings.HasPrefix(ltr, loc.FullAddress) {
		t.Errorf("buildLocationString(isRTL=false) should start with the address, got %q", ltr)
	}

	rtlOut := buildLocationString(loc, true)
	if !strings.HasPrefix(rtlOut, rlm+loc.FullAddress) {
		t.Errorf("buildLocationString(isRTL=true) should start with RLM + address, got %q", rtlOut)
	}
}

// TestSimpleTextualFallback_RTLMark applies the same check to the
// no-weather-data fallback path.
func TestSimpleTextualFallback_RTLMark(t *testing.T) {
	loc := &domain.Location{Name: "تلہ گنگ", TimeZone: "Asia/Karachi"}

	ltr := simpleTextualFallback(loc, false)
	if strings.Contains(ltr, rlm) {
		t.Errorf("simpleTextualFallback(isRTL=false) should not contain RLM, got %q", ltr)
	}

	rtlOut := simpleTextualFallback(loc, true)
	if !strings.Contains(rtlOut, rlm+loc.Name) {
		t.Errorf("simpleTextualFallback(isRTL=true) should contain RLM immediately before the location name, got %q", rtlOut)
	}
}
