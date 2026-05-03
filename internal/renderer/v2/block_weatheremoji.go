package v2

import (
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
)

// drawWeatherEmoji renders hourly weather emojis by reusing RenderConditionEmoji
// from the oneline package (via a minimal dummy context).
func drawWeatherEmoji(codes []int, opts *options.Options) string {
	if len(codes) == 0 {
		return "\n"
	}

	var b strings.Builder

	for _, code := range codes {
		ctx := &oneline.RenderContext{
			Data: &oneline.ParsedCurrentCondition{
				ConditionCode: strconv.Itoa(code),
			},
			Options: opts,
			// Location, Now, DataRaw can stay nil/zero — emoji renderer doesn't need them
		}

		b.WriteString(oneline.RenderConditionEmoji(ctx))
	}
	b.WriteRune('\n')
	return b.String()
}
