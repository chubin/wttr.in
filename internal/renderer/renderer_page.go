package renderer

import (
	"fmt"
	"strings"

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/localization"
)

// PageRenderer serves static embedded text pages (help, about, ...).
// It now fully supports multilanguage via the Localizer.
type PageRenderer struct{}

func NewPageRenderer() *PageRenderer {
	return &PageRenderer{}
}

// Render implements the Renderer interface.
func (p *PageRenderer) Render(query domain.Query, localizer localization.Localizer) (domain.RenderOutput, error) {
	loc := strings.TrimSpace(query.Options.Location)
	if loc == "" {
		loc = ":help"
	}

	// Normalize page name (:help → help, help → help)
	pageName := strings.TrimPrefix(loc, ":")
	if pageName == "" {
		pageName = "help"
	}

	// Use the new L10n wrapper for convenience
	l10n := localization.New(localizer, query.Options)

	// Try to load localized version first
	content, err := l10n.File(pageName + ".txt")
	if err == nil {
		return domain.RenderOutput{Content: []byte(content)}, nil
	}

	// Fallback: try English version via assets directly (backward compatibility + safety)
	filename := fmt.Sprintf("share/pages/%s.txt", pageName)

	contentBytes, err := assets.GetFile(filename)
	if err == nil {
		return domain.RenderOutput{Content: contentBytes}, nil
	}

	// Final fallback: generic error page (still localized!)
	// Should never happen.
	errorMsg := fmt.Sprintf("Page not found: %s\n", pageName)

	return domain.RenderOutput{
		Content: []byte(errorMsg),
	}, nil
}
