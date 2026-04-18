package renderer

import "github.com/chubin/wttr.in/internal/domain"

type V2Renderer struct{}

func (r *V2Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	// Stub: To be implemented
	return domain.RenderOutput{}, nil
}
