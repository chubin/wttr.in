package weather

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chubin/wttr.go/internal/location"
	"github.com/sirupsen/logrus"
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
func (l *cacheLocator) GetLocation(locationName string) (*Location, error) {
	// Use the existing Resolve method — it does cache lookup + upstream query + timezone
	// + caching of result
	raw, err := l.cache.Resolve(locationName)
	if err != nil {
		err1 := AppendToFile("/tmp/unknown-locations.txt", fmt.Sprintf("%s", locationName))
		if err1 != nil {
			logrus.Errorln(err1)
			return nil, err
		}

		return nil, err
	}

	// Convert between the two Location types
	lat, err := strconv.ParseFloat(raw.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude in cached location: %w", err)
	}

	lon, err := strconv.ParseFloat(raw.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude in cached location: %w", err)
	}

	// Country / CountryCode are not present in the existing Location type.
	// Options:
	//   1. Leave them empty (simplest)
	//   2. Parse them from Fullname (heuristic, fragile)
	//   3. Extend upstream query to return them (requires changing Nominatim parsers)
	//   4. Use a second service (e.g. reverse geocoding with country info)
	//
	// → for most use-cases option 1 is acceptable

	return &Location{
		Name: raw.Name, // already normalized
		// Country:      "",              // not available
		// CountryCode:  "",              // not available
		Latitude:    lat,
		Longitude:   lon,
		FullAddress: raw.Fullname,
	}, nil
}

// AppendToFile appends a string to the specified file with a timestamp
func AppendToFile(filename string, content string) error {
	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Get current timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Format the content with timestamp
	formattedContent := fmt.Sprintf("[%s] %s\n", timestamp, content)

	// Write to file
	_, err = file.WriteString(formattedContent)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
