package weather

import (
	"net/http"

	"github.com/chubin/wttr.go/internal/query"
)

func (s *WeatherService) upstreamRouting(
	opts *query.Options, r *http.Request, ipData *IPData, location *Location,
) (bool, *FormatOutput, error) {
	return false, nil, nil
}
