package location

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/ringsaturn/tzf"
)

var (
	tzFinder     tzf.F
	tzFinderOnce sync.Once
	tzFinderErr  error
)

func getTZFinder() (tzf.F, error) {
	tzFinderOnce.Do(func() {
		tzFinder, tzFinderErr = tzf.NewDefaultFinder()
	})
	return tzFinder, tzFinderErr
}

func enrichLocationWithTimezone(loc *Location) error {
	if loc == nil {
		return fmt.Errorf("nil location")
	}

	finder, err := getTZFinder()
	if err != nil {
		return err
	}

	lat, err := strconv.ParseFloat(loc.Lat, 64)
	if err != nil {
		return fmt.Errorf("invalid latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(loc.Lon, 64)
	if err != nil {
		return fmt.Errorf("invalid longitude: %w", err)
	}

	tzName := finder.GetTimezoneName(lon, lat)
	if tzName == "" {
		tzName = "UTC"
	}
	loc.Timezone = tzName
	return nil
}
