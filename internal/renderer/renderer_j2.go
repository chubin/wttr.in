package renderer

import (
	"encoding/json"

	"github.com/chubin/wttr.in/internal/domain"
)

// J2Renderer is a renderer that produces a minified JSON output of weather data
// with hourly entries removed for a more compact representation.
type J2Renderer struct{}

// Render processes the input Query and transforms the raw weather data into
// a minified JSON format with indentation. It removes all hourly entries from
// the weather data to reduce the output size.
//
// Parameters:
//   - query: A Query struct containing the weather data to be rendered.
//
// Returns:
//   - RenderOutput: A struct containing the rendered JSON data as bytes.
//   - error: An error if the JSON unmarshaling or marshaling process fails.
func (r *J2Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	// Check if weather data is available in the query
	if query.Weather == nil || len(*query.Weather) == 0 {
		return domain.RenderOutput{}, nil
	}

	// Define a temporary structure to hold the unmarshaled weather data
	var weatherData map[string]interface{}

	// Unmarshal the raw weather data into a generic map for manipulation
	err := json.Unmarshal(*query.Weather, &weatherData)
	if err != nil {
		return domain.RenderOutput{}, err
	}

	// Check if the "weather" key exists and contains forecast data
	if weather, ok := weatherData["weather"]; ok {
		if weatherSlice, ok := weather.([]interface{}); ok {
			// Iterate through each day's weather data and remove the "hourly" field
			for _, day := range weatherSlice {
				if dayMap, ok := day.(map[string]interface{}); ok {
					delete(dayMap, "hourly")
				}
			}
		}
	}

	// Marshal the modified weather data back to JSON with indentation for readability
	// Use a 2-space indent for a compact yet readable output
	minifiedJSON, err := json.MarshalIndent(weatherData, "", "  ")
	if err != nil {
		return domain.RenderOutput{}, err
	}

	// Return the rendered output as a RenderOutput struct
	return domain.RenderOutput{
		Content: minifiedJSON,
	}, nil
}
