package location

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// Check if the locationName is a GPS coordinate
	if isGPSCoordinate(locationName) {
		locationName, err = c.findNearestKnownLocation(locationName)
		if err != nil {
			log.Println("geo/location error:", locationName)
			return errorResponse(fmt.Sprint(err))
		}
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

// isGPSCoordinate checks if the given locationName is a GPS coordinate.
func isGPSCoordinate(locationName string) bool {
	parts := strings.Split(locationName, ",")
	if len(parts) != 2 {
		return false
	}

	_, err1 := strconv.ParseFloat(parts[0], 64)
	_, err2 := strconv.ParseFloat(parts[1], 64)

	return err1 == nil && err2 == nil
}

// findNearestKnownLocation finds the nearest known location for the given GPS coordinate.
func (c *Cache) findNearestKnownLocation(gps string) (string, error) {
	// Implement the logic to find the nearest known location for the given GPS coordinate.
	// This is a placeholder implementation and should be replaced with the actual logic.
	return "Montreal", nil
}
