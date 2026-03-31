package renderer

import (
	"encoding/json"
	"fmt"

	"github.com/chubin/wttr.in/internal/domain"
)

type J1Renderer struct{}

func (r *J1Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	if query.Weather == nil {
		return domain.RenderOutput{}, fmt.Errorf("no weather data provided")
	}

	var data struct {
		Data interface{} `json:"data"`
	}

	err := json.Unmarshal(*query.Weather, &data)
	if err != nil {
		return domain.RenderOutput{}, fmt.Errorf("invalid data format")
	}

	dataBytes, err := json.MarshalIndent(data.Data, "", "  ")
	if err != nil {
		return domain.RenderOutput{}, fmt.Errorf("invalid data format")
	}

	// Return the raw weather data as JSON bytes
	return domain.RenderOutput{
		Content: dataBytes,
	}, nil
}
