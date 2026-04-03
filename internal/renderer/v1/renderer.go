// Package v1 implements the original wttr.in "v1" rendering style
// (based on the classic wego-derived output from 2016).
package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
)

// V1Renderer renders weather in the classic wttr.in v1 style.
type V1Renderer struct {
	ansiEsc     *regexp.Regexp
	rightToLeft bool
}

// NewV1Renderer creates a new v1 renderer.
func NewV1Renderer() *V1Renderer {
	return &V1Renderer{
		ansiEsc: regexp.MustCompile(`\033\[.*?m`),
	}
}

// Render implements the weather.Renderer interface.
func (r *V1Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	if query.Weather == nil || len(*query.Weather) == 0 {
		return domain.RenderOutput{}, errors.New("no weather data provided")
	}

	// Unmarshal the raw weather data
	var data domain.Weather
	if err := json.Unmarshal(*query.Weather, &data); err != nil {
		return domain.RenderOutput{}, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}

	if len(data.CurrentCondition) == 0 {
		return domain.RenderOutput{}, errors.New("no current condition data available")
	}

	dataResp := ConvertWeather(data)

	opts := query.Options
	if opts == nil {
		opts = &options.Options{}
	}

	// Determine location name
	locationName := ""
	if query.Location != nil && query.Location.Name != "" {
		locationName = query.Location.Name
	} else if len(data.Request) > 0 && data.Request[0].Query != "" {
		locationName = data.Request[0].Query
	}
	if opts.Location != "" {
		locationName = opts.Location
	}

	// Right-to-left support
	r.rightToLeft = (opts.Lang == "he" || opts.Lang == "ar" || opts.Lang == "fa")

	// Build caption
	caption := "Weather report"
	if localized, ok := localizedCaption()[opts.Lang]; ok {
		caption = localized
	}

	var header string
	if r.rightToLeft {
		caption = locationName + " " + caption
		space := strings.Repeat(" ", 125-runewidth.StringWidth(caption))
		header = space + caption + "\n\n"
	} else {
		header = fmt.Sprintf("%s %s\n\n", caption, locationName)
	}

	// Current condition
	current := data.CurrentCondition[0]
	condLines := r.formatCond(convertCurrentConditionToCond(current), true, opts)

	// Build output using strings.Builder directly
	var sb strings.Builder
	sb.WriteString(header)

	// Write current condition with proper indentation
	for _, line := range condLines {
		if r.rightToLeft {
			sb.WriteString(strings.Repeat(" ", 94))
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Forecast days
	numDays := determineNumDays(opts)

	if numDays > 0 && len(data.Weather) > 0 {
		for i, day := range dataResp.Data.Weather {
			if i >= numDays {
				break
			}
			lines, err := r.printDay(day, opts)
			if err != nil {
				return domain.RenderOutput{}, err
			}
			for _, line := range lines {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}

	return domain.RenderOutput{
		Content: []byte(sb.String()),
	}, nil
}

// determineNumDays extracts forecast depth logic
func determineNumDays(opts *options.Options) int {
	if opts.Days > 0 {
		return opts.Days
	}
	if opts.CurrentOnly {
		return 0
	}
	if opts.CurrentPlusToday {
		return 1
	}
	if opts.CurrentPlusTwoDays {
		return 2
	}
	if opts.CurrentPlusThreeDays {
		return 3
	}
	return 3 // default
}

func convertCurrentConditionToCond(cc domain.CurrentCondition) cond {
	return cond{}
}
