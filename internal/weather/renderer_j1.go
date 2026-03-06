package weather

import "fmt"

type J1Renderer struct{}

func (r *J1Renderer) Render(query Query) (RenderOutput, error) {
	if query.Weather == nil {
		return RenderOutput{}, fmt.Errorf("no weather data provided")
	}

	// Return the raw weather data as JSON bytes
	return RenderOutput{
		Content: *query.Weather,
	}, nil
}
