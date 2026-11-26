package v2

import (
	"strings"
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/localization"
	"github.com/chubin/wttr.in/internal/options"
)

func TestTextualInformationUsesTranslatedCondition(t *testing.T) {
	weather := domain.WeatherRaw(`{
  "current_condition": [{
    "localObsDateTime": "2026-06-04 12:00",
    "weatherCode": "119",
    "weatherDesc": [{"value": "Cloudy"}],
    "lang_xx": [{"value": "Wolkig"}]
  }]
}`)
	l10n := localization.L10n{Text: func(key string) string { return key }}

	got := textualInformation(
		&domain.Query{Weather: &weather},
		&domain.Location{Name: "Berlin"},
		&options.Options{Lang: "de"},
		l10n,
	)

	if !strings.Contains(got, "Wolkig") || strings.Contains(got, "Cloudy") {
		t.Fatalf("textualInformation() does not use the translated condition:\n%s", got)
	}
}
