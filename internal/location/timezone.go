package location

import (
	"fmt"
	"strconv"

	"github.com/ringsaturn/tzf"
)

// LatLonToTimezone returns the IANA timezone name for the given coordinates.
// Falls back to "UTC" if the lookup fails or returns an empty result.
func LatLonToTimezone(lat, lon float64) (string, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return "UTC", err
	}
	tzName := finder.GetTimezoneName(lon, lat)
	if tzName == "" {
		tzName = "UTC"
	}
	return tzName, nil
}

func enrichLocationWithTimezone(loc *Location) error {
	if loc == nil {
		return fmt.Errorf("nil location")
	}

	// Convert between the two Location types
	lat, err := strconv.ParseFloat(loc.Lat, 64)
	if err != nil {
		return fmt.Errorf("invalid latitude in cached location: %w", err)
	}

	lon, err := strconv.ParseFloat(loc.Lon, 64)
	if err != nil {
		return fmt.Errorf("invalid longitude in cached location: %w", err)
	}

	tzName, err := LatLonToTimezone(lat, lon)
	if err != nil {
		return err
	}

	loc.Timezone = tzName

	return nil
}
