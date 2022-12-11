package location

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
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
		db, err = godb.Open(sqlite.Adapter, config.Geo.IPCacheDB)
		if err != nil {
			return nil, err
		}

		// Needed for "upsert" implementation in Put()
		db.UseErrorParser()
	}

	return &Cache{
		config:        config,
		db:            db,
		indexField:    "name",
		filesCacheDir: config.Geo.LocationCache,
	}, nil
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
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Cache) Put(addr string, loc *Location) error {
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
