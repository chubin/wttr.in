package location

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chubin/wttr.in/internal/routing"
)

// Response provides routing interface to the geo cache.
func (c *Cache) Response(r *http.Request) *routing.Cadre {
	var (
		locationName = r.URL.Query().Get("location")
		loc          *Location
		bytes        []byte
		err          error
	)

	if locationName == "" {
		return errorResponse("location is not specified")
	}

	if strings.Contains(locationName, ".html") || strings.Contains(locationName, ".txt") {
		return errorResponse("invalid location")
	}

	loc, err = c.Resolve(locationName)
	if err != nil {
		log.Println("geo/location error:", locationName, r.RemoteAddr)

		return errorResponse(fmt.Sprint(err))
	}

	bytes, err = json.Marshal(loc)
	if err != nil {
		return errorResponse(fmt.Sprint(err))
	}

	return &routing.Cadre{Body: bytes}
}

func errorResponse(s string) *routing.Cadre {
	return &routing.Cadre{Body: []byte(
		fmt.Sprintf(`{"error": %q}`, s),
	)}
}
