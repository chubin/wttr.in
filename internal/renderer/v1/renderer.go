// Package v1 implements the original wttr.in "v1" rendering style
// (based on the classic wego-derived output from 2016).
package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
	// if you have common helpers
)

// V1Renderer renders weather in the classic wttr.in v1 style.
type V1Renderer struct {
	ansiEsc *regexp.Regexp
}

// NewV1Renderer creates a new v1 renderer.
func NewV1Renderer() *V1Renderer {
	return &V1Renderer{
		ansiEsc: regexp.MustCompile(`\033.*?m`),
	}
}

// Render implements the weather.Renderer interface.
func (r *V1Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	if query.Weather == nil || len(*query.Weather) == 0 {
		return domain.RenderOutput{}, errors.New("no weather data provided")
	}

	// Unmarshal the raw weather data into the domain model
	var data domain.Weather
	if err := json.Unmarshal(*query.Weather, &data); err != nil {
		return domain.RenderOutput{}, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}

	if len(data.CurrentCondition) == 0 {
		return domain.RenderOutput{}, errors.New("no current condition data available")
	}

	opts := query.Options
	if opts == nil {
		opts = &options.Options{} // fallback to zero value
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

	// Right-to-left support for certain languages
	rightToLeft := opts.Lang == "he" || opts.Lang == "ar" || opts.Lang == "fa"

	// Build caption
	caption := "Weather report"
	if localized, ok := localizedCaption()[opts.Lang]; ok {
		caption = localized
	}

	var header string
	if rightToLeft {
		caption = locationName + " " + caption
		space := strings.Repeat(" ", 125-runewidth.StringWidth(caption))
		header = space + caption + "\n\n"
	} else {
		header = fmt.Sprintf("%s %s\n\n", caption, locationName)
	}

	// Current condition
	current := data.CurrentCondition[0]
	condLines := r.formatCond(make([]string, 5), current, true)

	// Build output
	var sb strings.Builder
	sb.WriteString(header)

	stdout := colorable.NewColorable(&sb) // simulate colored output into buffer

	for _, line := range condLines {
		if rightToLeft {
			fmt.Fprint(stdout, strings.Repeat(" ", 94))
		} else {
			fmt.Fprint(stdout, " ")
		}
		fmt.Fprintln(stdout, line)
	}

	// Forecast days
	numDays := 3
	if opts.Days > 0 {
		numDays = opts.Days
	} else if opts.CurrentOnly {
		numDays = 0
	} else if opts.CurrentPlusToday {
		numDays = 1
	} else if opts.CurrentPlusTwoDays {
		numDays = 2
	} else if opts.CurrentPlusThreeDays {
		numDays = 3
	}

	if numDays > 0 && len(data.Weather) > 0 {
		for i, day := range data.Weather {
			if i >= numDays {
				break
			}
			lines, err := r.printDay(day)
			if err != nil {
				return domain.RenderOutput{}, err
			}
			for _, line := range lines {
				fmt.Fprintln(stdout, line)
			}
		}
	}

	return domain.RenderOutput{
		Content: sb.Bytes(),
	}, nil
}
