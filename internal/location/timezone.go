package location

import (
	"fmt"
	"strconv"

	"github.com/ringsaturn/tzf"
)

func enrichLocationWithTimezone(loc *Location) error {
	if loc == nil {
		return fmt.Errorf("nil location")
	}

	finder, err := tzf.NewDefaultFinder() // loads embedded timezone data
	if err != nil {
		return err
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

	tzName := finder.GetTimezoneName(lat, lon)
	if tzName == "" {
		tzName = "UTC" // fallback
	}

	loc.Timezone = tzName // assuming you added TimeZone string to Location struct

	return nil
}
