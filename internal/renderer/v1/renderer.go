// Package v1 implements the original wttr.in "v1" rendering style
// (based on the classic wego-derived output from 2016).
package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/clipperhouse/displaywidth"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/localization"
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
func (r *V1Renderer) Render(query domain.Query, localizer localization.Localizer) (domain.RenderOutput, error) {
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

	l10n := localization.New(localizer, query.Options)
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
	r.rightToLeft = l10n.IsRTL()

	// Build caption
	caption := l10n.Text("CAPTION_WEATHER_REPORT_FOR")

	var header string
	if opts.Quiet || opts.NoCaption {
		header = locationName + "\n\n"
	} else {
		if r.rightToLeft {
			caption = locationName + " " + caption
			padding := 125 - displaywidth.String(caption)
			if padding < 0 {
				padding = 0
			}
			space := strings.Repeat(" ", padding)
			header = space + caption + "\n\n"
		} else {
			header = fmt.Sprintf("%s %s\n\n", caption, locationName)
		}
	}

	// Current condition
	current := data.CurrentCondition[0]
	cond := convertCurrentConditionToCond(current)
	condLines := r.formatCond("", cond, true, opts)

	// Build output using strings.Builder directly
	var sb strings.Builder
	if !opts.Superquiet && !opts.NoCity {
		sb.WriteString(header)
	}

	// Write current condition with proper indentation
	for _, line := range condLines {
		if r.rightToLeft {
			sb.WriteString(strings.Repeat(" ", 94))
		} else {
			sb.WriteString("")
		}
		sb.WriteString(strings.TrimRight(line, " "))
		sb.WriteString("\n")
	}

	// Forecast days
	numDays := determineNumDays(opts)

	if numDays > 0 && len(data.Weather) > 0 {
		for i, day := range dataResp.Data.Weather {
			if i >= numDays {
				break
			}
			lines, err := r.printDay(day, opts, l10n)
			if err != nil {
				return domain.RenderOutput{}, err
			}
			for _, line := range lines {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}

	// Location line:
	if !opts.CurrentOnly {

		if !opts.Quiet && !opts.Superquiet && !opts.NoCity {
			fullAddress := query.Location.FullAddress
			if r.rightToLeft {
				// Prepend a right-to-left mark so terminals render the
				// location name in the correct visual order instead of
				// treating it as left-to-right text (issue #932).
				fullAddress = rlm + fullAddress
			}
			sb.WriteString(fmt.Sprintf("%s: %s [%v,%v]\n",
				l10n.Text("LOCATION"),
				fullAddress, query.Location.Latitude, query.Location.Longitude))
		}

		if opts.Output != "html" && !opts.NoFollowLine {
			followICforUpdates := l10n.Text("FOLLOW_ME")
			sb.WriteString("\n" + followICforUpdates + "\n")
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

// convertCurrentConditionToCond converts domain.CurrentCondition to v1.cond
func convertCurrentConditionToCond(cc domain.CurrentCondition) cond {
	c := cond{
		FeelsLikeC:     parseInt(cc.FeelsLikeC),
		PrecipMM:       parseFloat32(cc.PrecipMM),
		TempC2:         parseInt(cc.TempC),
		VisibleDistKM:  parseInt(cc.Visibility),
		WeatherCode:    parseInt(cc.WeatherCode),
		WindspeedKmph:  parseInt(cc.WindspeedKmph),
		Winddir16Point: cc.Winddir16Point,
		WeatherDesc: []struct{ Value string }{
			{Value: getTranslatedDesc(cc.WeatherDesc, cc.LangXX, "Unknown")},
		},
	}

	return c
}
