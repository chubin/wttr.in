package weather

import (
	"strings"
	"testing"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/formatter"
)

func TestSelectFormatterJSONContentType(t *testing.T) {
	s := &WeatherService{
		FormatterMap: map[string]Formatter{
			"text": &formatter.TextFormatter{},
			"json": &formatter.JSONFormatter{},
		},
	}

	out, err := s.selectFormatter("json").Format(nil, &domain.RenderOutput{Content: []byte(`{"ok":true}`)})
	if err != nil {
		t.Fatalf("json formatter: %v", err)
	}
	if !strings.HasPrefix(out.ContentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", out.ContentType)
	}

	textOut, err := s.selectFormatter("text").Format(nil, &domain.RenderOutput{Content: []byte("ok")})
	if err != nil {
		t.Fatalf("text formatter: %v", err)
	}
	if !strings.HasPrefix(textOut.ContentType, "text/plain") {
		t.Fatalf("Content-Type = %q, want text/plain", textOut.ContentType)
	}
}
