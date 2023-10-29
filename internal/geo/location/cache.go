package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
	log "github.com/sirupsen/logrus"
	"github.com/zsefvlol/timezonemapper"

	"github.com/chubin/wttr.in/internal/config"
	"github.com/chubin/wttr.in/internal/types"
)

// Cache is an implemenation of DB/file-based cache.
//
// At the moment, it is an implementation for the location cache,
// but it should be generalized to cache everything.
type Cache struct {
	config        *config.Config
	db            *godb.DB
	searcher      *Searcher
	indexField    string
	filesCacheDir string
}

// NewCache returns new cache reader for the specified config.
func NewCache(config *config.Config) (*Cache, error) {
	var (
		db  *godb.DB
		err error
	)

	if config.Geo.LocationCacheType == types.CacheTypeDB {
		log.Debugln("using db for location cache")
		db, err = godb.Open(sqlite.Adapter, config.Geo.LocationCacheDB)
		if err != nil {
			return nil, err
		}

		log.Debugln("db file:", config.Geo.LocationCacheDB)

		// Needed for "upsert" implementation in Put()
		db.UseErrorParser()
	}

	return &Cache{
		config:        config,
		db:            db,
		indexField:    "name",
		filesCacheDir: config.Geo.LocationCache,
		searcher:      NewSearcher(config),
	}, nil
}

// Resolve returns location information for specified location.
// If the information is found in the cache, it is returned.
// If it is not found, the external service is queried,
// and the result is stored in the cache.
func (c *Cache) Resolve(location string) (*Location, error) {
	location = normalizeLocationName(location)

	loc, err := c.Read(location)
	if !errors.Is(err, types.ErrNotFound) {
		return loc, err
	}

	log.Debugln("geo/location: not found in cache:", location)
	loc, err = c.searcher.Search(location)
	if err != nil {
		return nil, err
	}

	loc.Name = location
	loc.Timezone = latLngToTimezoneString(loc.Lat, loc.Lon)

	err = c.Put(location, loc)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// Read returns location information from the cache, if found,
// or types.ErrNotFound if not found. If the entry is found, but its format
// is invalid, types.ErrInvalidCacheEntry is returned.
func (c *Cache) Read(addr string) (*Location, error) {
	if c.config.Geo.LocationCacheType == types.CacheTypeFiles {
		return c.readFromCacheFile(addr)
	}

	return c.readFromCacheDB(addr)
}

func (c *Cache) readFromCacheFile(name string) (*Location, error) {
	var (
		fileLoc = struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Timezone  string  `json:"timezone"`
			Address   string  `json:"address"`
		}{}
		location Location
	)

	bytes, err := os.ReadFile(c.cacheFile(name))
	if err != nil {
		return nil, types.ErrNotFound
	}
	err = json.Unmarshal(bytes, &fileLoc)
	if err != nil {
		return nil, err
	}

	// normalize name
	name = strings.TrimSpace(
		strings.TrimRight(
			strings.TrimLeft(name, `"`), `"`))

	timezone := fileLoc.Timezone
	if timezone == "" {
		timezone = timezonemapper.LatLngToTimezoneString(fileLoc.Latitude, fileLoc.Longitude)
	}

	location = Location{
		Name:     name,
		Lat:      fmt.Sprint(fileLoc.Latitude),
		Lon:      fmt.Sprint(fileLoc.Longitude),
		Timezone: timezone,
		Fullname: fileLoc.Address,
	}

	return &location, nil
}

func (c *Cache) readFromCacheDB(addr string) (*Location, error) {
	result := Location{}
	err := c.db.Select(&result).
		Where(c.indexField+" = ?", addr).
		Do()

	if strings.Contains(fmt.Sprint(err), "no rows in result set") {
		return nil, types.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("readFromCacheDB: %w", err)
	}

	return &result, nil
}

func (c *Cache) Put(addr string, loc *Location) error {
	log.Infoln("geo/location: storing in cache:", loc)
	if c.config.Geo.IPCacheType == types.CacheTypeDB {
		return c.putToCacheDB(loc)
	}

	return c.putToCacheFile(addr, loc)
}

func (c *Cache) putToCacheDB(loc *Location) error {
	err := c.db.Insert(loc).Do()
	if strings.Contains(fmt.Sprint(err), "UNIQUE constraint failed") {
		return c.db.Update(loc).Do()
	}

	return err
}

func (c *Cache) putToCacheFile(addr string, loc fmt.Stringer) error {
	return os.WriteFile(c.cacheFile(addr), []byte(loc.String()), 0o600)
}

// cacheFile returns path to the cache entry for addr.
func (c *Cache) cacheFile(item string) string {
	return path.Join(c.filesCacheDir, item)
}

// normalizeLocationName converts name into the standard location form
// with the following steps:
// - remove excessive spaces,
// - remove quotes,
// - convert to lover case.
func normalizeLocationName(name string) string {
	name = strings.ReplaceAll(name, `"`, " ")
	name = strings.ReplaceAll(name, `'`, " ")
	name = strings.TrimSpace(name)
	name = strings.Join(strings.Fields(name), " ")

	return strings.ToLower(name)
}

// latLngToTimezoneString returns timezone for lat, lon,
// or an empty string if they are invalid.
func latLngToTimezoneString(lat, lon string) string {
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Errorln("geoloc: latLngToTimezoneString:", err)

		return ""
	}
	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Errorln("geoloc: latLngToTimezoneString:", err)

		return ""
	}

	return timezonemapper.LatLngToTimezoneString(latFloat, lonFloat)
}
