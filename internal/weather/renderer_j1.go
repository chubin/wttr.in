package weather

import (
	"encoding/json"
	"fmt"
)

type J1Renderer struct{}

func (r *J1Renderer) Render(query Query) (RenderOutput, error) {
	if query.Weather == nil {
		return RenderOutput{}, fmt.Errorf("no weather data provided")
	}

	var data struct {
		Data interface{} `json:"data"`
	}

	err := json.Unmarshal(*query.Weather, &data)
	if err != nil {
		return RenderOutput{}, fmt.Errorf("invalid data format")
	}

	dataBytes, err := json.MarshalIndent(data.Data, "", "  ")
	if err != nil {
		return RenderOutput{}, fmt.Errorf("invalid data format")
	}

	// Return the raw weather data as JSON bytes
	return RenderOutput{
		Content: dataBytes,
	}, nil
}
