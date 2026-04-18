package renderer

import (
	"fmt"

	"github.com/chubin/wttr.in/internal/domain"
)

type J1Renderer struct{}

func (r *J1Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	if query.Weather == nil {
		return domain.RenderOutput{}, fmt.Errorf("no weather data provided")
	}

	// Return the raw weather data as JSON bytes
	return domain.RenderOutput{
		Content: *query.Weather,
	}, nil
}
