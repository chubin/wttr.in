package ip

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/config"
)

var (
	ErrNotFound          = errors.New("cache entry not found")
	ErrInvalidCacheEntry = errors.New("invalid cache entry format")
)

// Location information.
type Location struct {
	CountryCode string
	Country     string
	Region      string
	City        string
	Latitude    float64
	Longitude   float64
}

// Cache provides access to the IP Geodata cache.
type Cache struct {
	config *config.Config
}

// NewCache returns new cache reader for the specified config.
func NewCache(config *config.Config) *Cache {
	return &Cache{
		config: config,
	}
}

// Read returns location information from the cache, if found,
// or ErrNotFound if not found. If the entry is found, but its format
// is invalid, ErrInvalidCacheEntry is returned.
//
// Format:
//
//     [CountryCode];Country;Region;City;[Latitude];[Longitude]
//
// Example:
//
//     DE;Germany;Free and Hanseatic City of Hamburg;Hamburg;53.5736;9.9782
//
func (c *Cache) Read(addr string) (*Location, error) {
	bytes, err := os.ReadFile(c.cacheFile(addr))
	if err != nil {
		return nil, ErrNotFound
	}
	return parseCacheEntry(string(bytes))
}

// cacheFile retuns path to the cache entry for addr.
func (c *Cache) cacheFile(addr string) string {
	return path.Join(c.config.Geo.IPCache, addr)
}

// parseCacheEntry parses the location cache entry s,
// and return location, or error, if the cache entry is invalid.
func parseCacheEntry(s string) (*Location, error) {
	var (
		lat  float64 = -1000
		long float64 = -1000
		err  error
	)

	parts := strings.Split(s, ";")
	if len(parts) < 4 {
		return nil, ErrInvalidCacheEntry
	}

	if len(parts) >= 6 {
		lat, err = strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return nil, ErrInvalidCacheEntry
		}

		long, err = strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, ErrInvalidCacheEntry
		}
	}

	return &Location{
		CountryCode: parts[0],
		Country:     parts[1],
		Region:      parts[2],
		City:        parts[3],
		Latitude:    lat,
		Longitude:   long,
	}, nil
}
