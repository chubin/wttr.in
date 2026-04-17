package renderer

import (
	"fmt"
	"strings"

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/chubin/wttr.in/internal/domain"
)

// PageRenderer serves static embedded text pages.
// It is used when the view is "page" or when the location starts with ":".
type PageRenderer struct{}

func NewPageRenderer() *PageRenderer {
	return &PageRenderer{}
}

// Render implements the Renderer interface defined in the domain layer.
func (p *PageRenderer) Render(query domain.Query) (domain.RenderOutput, error) {
	// Extract location (which contains the page name when using :page syntax)
	loc := strings.TrimSpace(query.Options.Location)
	if loc == "" {
		loc = ":help"
	}

	// Normalize: support both ":help" and "help"
	pageName := strings.TrimPrefix(loc, ":")
	if pageName == "" {
		pageName = "help"
	}

	// Derive embedded filename: share/pages/{NAME}.txt
	filename := fmt.Sprintf("share/pages/%s.txt", pageName)

	content, err := assets.GetFile(filename)
	if err != nil {
		// File not found → return a helpful error page (still as RenderOutput)
		errorMsg := fmt.Sprintf("Page not found: :%s\n", pageName)
		return domain.RenderOutput{
			Content: []byte(errorMsg),
		}, nil
	}

	return domain.RenderOutput{
		Content: content,
	}, nil
}
