// Package location provides geocoding cache functionality with file- or DB-based storage.
package location

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
	log "github.com/sirupsen/logrus"
	"github.com/zsefvlol/timezonemapper"

	"github.com/chubin/wttr.go/internal/types"
)

var (
	CacheTypeDB    = "db"
	CacheTypeFiles = "files"

	// Tuning constants – adjust according to workload & memory budget
	defaultLRUSize         = 8192
	defaultMaxOpenConns    = 16
	defaultMaxIdleConns    = 4
	defaultConnMaxIdleTime = 5 * time.Minute
)

type Config struct {
	LocationCacheType string `yaml:"locationCacheType"`
	LocationCacheDB   string `yaml:"locationCacheDb"`
	LocationCache     string `yaml:"locationCache"`
	NominatimServers  []struct {
		Name  string `yaml:"name"`
		Type  string `yaml:"type"`
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"nominatim"`
}

// Location represents a cached geocoding result.
type Location struct {
	Name     string `db:"name, key" json:"name"`
	Lat      string `db:"lat" json:"latitude"`
	Lon      string `db:"lon" json:"longitude"`
	Timezone string `db:"timezone" json:"timezone"`
	Fullname string `db:"displayName" json:"address"`
}

// String returns string representation of location.
func (l *Location) String() string {
	bytes, err := json.Marshal(l)
	if err != nil {
		// should never happen
		log.Fatalln(err)
	}

	return string(bytes)
}

// Cache combines in-memory LRU, SQLite read pool, and write support.
type Cache struct {
	config *Config

	// ── Persistent storage ───────────────────────────────
	dbPath      string
	godbWrite   *godb.DB // used for writes (Insert/Update)
	sqlReadPool *sql.DB  // read-only connection pool

	// ── Fast in-memory cache ─────────────────────────────
	lru   *lru.Cache[string, *Location]
	lruMu sync.Mutex // only for lazy init

	// ── Synchronization ──────────────────────────────────
	writeMu sync.Mutex

	searcher      *Searcher
	indexField    string
	filesCacheDir string
}

// NewCache creates a new cache instance according to the configuration.
func NewCache(config *Config) (*Cache, error) {
	c := &Cache{
		config:        config,
		indexField:    "name",
		filesCacheDir: config.LocationCache,
		searcher:      NewSearcher(config),
	}

	if config.LocationCacheType == CacheTypeDB {
		log.Debugln("using db for location cache")
		c.dbPath = config.LocationCacheDB

		// Write connection via godb
		db, err := godb.Open(sqlite.Adapter, c.dbPath)
		if err != nil {
			return nil, fmt.Errorf("cannot open godb write connection: %w", err)
		}
		db.UseErrorParser()
		c.godbWrite = db

		// Read pool via database/sql (better concurrency control)
		sqlDB, err := sql.Open("sqlite3", c.dbPath+"?_journal_mode=WAL&_busy_timeout=5000&mode=ro")
		if err != nil {
			return nil, fmt.Errorf("cannot open sqlite read pool: %w", err)
		}

		sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
		sqlDB.SetMaxIdleConns(defaultMaxIdleConns)
		sqlDB.SetConnMaxIdleTime(defaultConnMaxIdleTime)

		c.sqlReadPool = sqlDB

		log.Debugf("SQLite read pool initialized (maxOpen=%d, maxIdle=%d)",
			defaultMaxOpenConns, defaultMaxIdleConns)
	}

	// LRU is initialized lazily on first use
	c.lru = nil

	// Ensure DB schema exists (if using DB mode)
	if err := c.InitDB(false); err != nil {
		return nil, err
	}

	return c, nil
}

// Resolve returns location data — first from cache (LRU → persistent), then from external search.
func (c *Cache) Resolve(location string) (*Location, error) {
	defer log.Debugln("Resolve() finished")

	location = normalizeLocationName(location)

	// 1. Fastest path: LRU
	if loc := c.getFromLRU(location); loc != nil {
		return loc, nil
	}

	// 2. Persistent cache
	loc, err := c.Read(location)
	if err == nil {
		c.putToLRU(location, loc) // warm the LRU
		return loc, nil
	}
	if !errors.Is(err, types.ErrNotFound) {
		return nil, err
	}

	// 3. External lookup
	loc, err = c.searcher.Search(location)
	if err != nil {
		return nil, err
	}

	loc.Name = location
	loc.Timezone = latLngToTimezoneString(loc.Lat, loc.Lon)

	// 4. Store result
	if err = c.Put(location, loc); err != nil {
		log.Warnln("failed to persist location to cache:", err)
		// still return fresh data
	}

	c.putToLRU(location, loc)

	return loc, nil
}

// Read fetches location from cache (LRU already checked in Resolve).
func (c *Cache) Read(addr string) (*Location, error) {
	addr = normalizeLocationName(addr)

	if c.config.LocationCacheType == CacheTypeFiles {
		return c.readFromCacheFile(addr)
	}
	return c.readFromCacheDB(addr)
}

func (c *Cache) readFromCacheFile(name string) (*Location, error) {
	log.Debugln("readFromCacheFile started")
	defer log.Debugln("readFromCacheFile finished")

	data, err := os.ReadFile(c.cacheFile(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var fileLoc struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
		Address   string  `json:"address"`
	}

	if err = json.Unmarshal(data, &fileLoc); err != nil {
		return nil, fmt.Errorf("invalid cache file format: %w", err)
	}

	tz := fileLoc.Timezone
	if tz == "" {
		tz = timezonemapper.LatLngToTimezoneString(fileLoc.Latitude, fileLoc.Longitude)
	}

	return &Location{
		Name:     name,
		Lat:      fmt.Sprintf("%.6f", fileLoc.Latitude),
		Lon:      fmt.Sprintf("%.6f", fileLoc.Longitude),
		Timezone: tz,
		Fullname: fileLoc.Address,
	}, nil
}

func (c *Cache) readFromCacheDB(addr string) (*Location, error) {
	start := time.Now()
	defer func() {
		log.Debugf("readFromCacheDB finished (time: %v)", time.Since(start))
	}()

	var loc Location
	err := c.sqlReadPool.QueryRow(
		`SELECT name, lat, lon, timezone, displayName
		 FROM Location
		 WHERE name = ? LIMIT 1`,
		addr,
	).Scan(
		&loc.Name, &loc.Lat, &loc.Lon, &loc.Timezone, &loc.Fullname,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite read failed: %w", err)
	}

	return &loc, nil
}

// Put stores the location in the configured backend.
func (c *Cache) Put(addr string, loc *Location) error {
	log.Infoln("geo/location: storing in cache:", loc)

	addr = normalizeLocationName(addr)

	if c.config.LocationCacheType == CacheTypeDB {
		return c.putToCacheDB(loc)
	}
	return c.putToCacheFile(addr, loc)
}

func (c *Cache) putToCacheDB(loc *Location) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	log.Debugln("putToCacheDB started")
	defer log.Debugln("putToCacheDB finished")

	err := c.godbWrite.Insert(loc).Do()
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique constraint") {
		return c.godbWrite.Update(loc).Do()
	}
	return err
}

func (c *Cache) putToCacheFile(addr string, loc *Location) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	log.Debugln("putToCacheFile started")
	defer log.Debugln("putToCacheFile finished")

	// For compatibility with old file format
	data := map[string]interface{}{
		"latitude":  mustParseFloat(loc.Lat),
		"longitude": mustParseFloat(loc.Lon),
		"timezone":  loc.Timezone,
		"address":   loc.Fullname,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.cacheFile(addr), bytes, 0o600)
}

// ── LRU helpers ────────────────────────────────────────────────────────────────

func (c *Cache) initLRUOnce() {
	c.lruMu.Lock()
	defer c.lruMu.Unlock()
	if c.lru != nil {
		return
	}

	l, err := lru.New[string, *Location](defaultLRUSize)
	if err != nil {
		log.Errorln("failed to initialize LRU cache:", err)
		return
	}

	c.lru = l
	log.Infof("LRU cache initialized (size = %d)", defaultLRUSize)
}

func (c *Cache) getFromLRU(key string) *Location {
	if c.lru == nil {
		c.initLRUOnce()
		if c.lru == nil {
			return nil
		}
	}
	val, ok := c.lru.Get(key)
	if ok {
		return val
	}
	return nil
}

func (c *Cache) putToLRU(key string, val *Location) {
	if c.lru == nil {
		c.initLRUOnce()
		if c.lru == nil {
			return
		}
	}
	c.lru.Add(key, val)
}

// ── Utility functions ──────────────────────────────────────────────────────────

func (c *Cache) cacheFile(item string) string {
	return path.Join(c.filesCacheDir, item)
}

func normalizeLocationName(name string) string {
	name = strings.ReplaceAll(name, `"`, " ")
	name = strings.ReplaceAll(name, `'`, " ")
	name = strings.TrimSpace(name)
	name = strings.Join(strings.Fields(name), " ")
	return strings.ToLower(name)
}

func latLngToTimezoneString(lat, lon string) string {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Errorln("latLngToTimezoneString: invalid lat:", err)
		return ""
	}
	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Errorln("latLngToTimezoneString: invalid lon:", err)
		return ""
	}
	return timezonemapper.LatLngToTimezoneString(latF, lonF)
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// Close cleans up resources (call on shutdown).
func (c *Cache) Close() error {
	var errs []error

	if c.sqlReadPool != nil {
		if err := c.sqlReadPool.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.godbWrite != nil {
		if err := c.godbWrite.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}
