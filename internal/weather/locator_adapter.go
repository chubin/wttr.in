package weather

import (
	"fmt"
	"strconv"

	"github.com/chubin/wttr.go/internal/location"
)

// cacheLocator implements Locator using wttr.in's Cache + Resolve logic
type cacheLocator struct {
	cache *location.Cache
}

// NewCacheLocator creates a Locator backed by wttr.in-style location cache
func NewCacheLocator(cache *location.Cache) Locator {
	if cache == nil {
		panic("cache must not be nil")
	}
	return &cacheLocator{cache: cache}
}

// GetLocation implements Locator interface
func (l *cacheLocator) GetLocation(locationName string) (Location, error) {
	// Use the existing Resolve method — it does cache lookup + upstream query + timezone
	// + caching of result
	raw, err := l.cache.Resolve(locationName)
	if err != nil {
		return Location{}, err
	}

	// Convert between the two Location types
	lat, err := strconv.ParseFloat(raw.Lat, 64)
	if err != nil {
		return Location{}, fmt.Errorf("invalid latitude in cached location: %w", err)
	}

	lon, err := strconv.ParseFloat(raw.Lon, 64)
	if err != nil {
		return Location{}, fmt.Errorf("invalid longitude in cached location: %w", err)
	}

	// Country / CountryCode are not present in the existing Location type.
	// Options:
	//   1. Leave them empty (simplest)
	//   2. Parse them from Fullname (heuristic, fragile)
	//   3. Extend upstream query to return them (requires changing Nominatim parsers)
	//   4. Use a second service (e.g. reverse geocoding with country info)
	//
	// → for most use-cases option 1 is acceptable

	return Location{
		Name: raw.Name, // already normalized
		// Country:      "",              // not available
		// CountryCode:  "",              // not available
		Latitude:    lat,
		Longitude:   lon,
		FullAddress: raw.Fullname,
	}, nil
}
