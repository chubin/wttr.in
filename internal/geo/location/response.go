package location

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	loc, err = c.Resolve(locationName)
	if err != nil {
		log.Println("geo/location error:", locationName)

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
