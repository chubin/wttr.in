package ip

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/config"
	"github.com/chubin/wttr.in/internal/routing"
	"github.com/chubin/wttr.in/internal/types"
	"github.com/chubin/wttr.in/internal/util"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

var (
	ErrNotFound          = errors.New("cache entry not found")
	ErrInvalidCacheEntry = errors.New("invalid cache entry format")
)

// Location information.
type Location struct {
	IP          string  `db:"ip,key"`
	CountryCode string  `db:"countryCode"`
	Country     string  `db:"country"`
	Region      string  `db:"region"`
	City        string  `db:"city"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
}

func (l *Location) String() string {
	if l.Latitude == -1000 {
		return fmt.Sprintf(
			"%s;%s;%s;%s",
			l.CountryCode, l.CountryCode, l.Region, l.City)
	}
	return fmt.Sprintf(
		"%s;%s;%s;%s;%v;%v",
		l.CountryCode, l.CountryCode, l.Region, l.City, l.Latitude, l.Longitude)
}

// Cache provides access to the IP Geodata cache.
type Cache struct {
	config *config.Config
	db     *godb.DB
}

// NewCache returns new cache reader for the specified config.
func NewCache(config *config.Config) (*Cache, error) {
	db, err := godb.Open(sqlite.Adapter, config.Geo.IPCacheDB)
	if err != nil {
		return nil, err
	}

	// Needed for "upsert" implementation in Put()
	db.UseErrorParser()

	return &Cache{
		config: config,
		db:     db,
	}, nil
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
	if c.config.Geo.CacheType == types.CacheTypeDB {
		return c.readFromCacheDB(addr)
	}
	return c.readFromCacheFile(addr)
}

func (c *Cache) readFromCacheFile(addr string) (*Location, error) {
	bytes, err := os.ReadFile(c.cacheFile(addr))
	if err != nil {
		return nil, ErrNotFound
	}
	return NewLocationFromString(addr, string(bytes))
}

func (c *Cache) readFromCacheDB(addr string) (*Location, error) {
	result := Location{}
	err := c.db.Select(&result).
		Where("IP = ?", addr).
		Do()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Cache) Put(addr string, loc *Location) error {
	if c.config.Geo.CacheType == types.CacheTypeDB {
		return c.putToCacheDB(addr, loc)
	}
	return c.putToCacheFile(addr, loc)
}

func (c *Cache) putToCacheDB(addr string, loc *Location) error {
	err := c.db.Insert(loc).Do()
	// it should work like this:
	//
	//   target := dberror.UniqueConstraint{}
	//   if errors.As(err, &target) {
	//
	// See: https://github.com/samonzeweb/godb/pull/23
	//
	// But for some reason it does not work,
	// so the dirty hack is used:
	if strings.Contains(fmt.Sprint(err), "UNIQUE constraint failed") {
		return c.db.Update(loc).Do()
	}
	return err
}

func (c *Cache) putToCacheFile(addr string, loc *Location) error {
	return os.WriteFile(c.cacheFile(addr), []byte(loc.String()), 0644)
}

// cacheFile retuns path to the cache entry for addr.
func (c *Cache) cacheFile(addr string) string {
	return path.Join(c.config.Geo.IPCache, addr)
}

// NewLocationFromString parses the location cache entry s,
// and return location, or error, if the cache entry is invalid.
func NewLocationFromString(addr, s string) (*Location, error) {
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
		IP:          addr,
		CountryCode: parts[0],
		Country:     parts[1],
		Region:      parts[2],
		City:        parts[3],
		Latitude:    lat,
		Longitude:   long,
	}, nil
}

// Reponse provides routing interface to the geo cache.
//
// Temporary workaround to switch IP addresses handling to the Go server.
// Handles two queries:
//
//		/:geo-ip-put?ip=IP&value=VALUE
//      /:geo-ip-get?ip=IP
//
func (c *Cache) Response(r *http.Request) *routing.Cadre {
	var (
		respERR = &routing.Cadre{Body: []byte("ERR")}
		respOK  = &routing.Cadre{Body: []byte("OK")}
	)

	if ip := util.ReadUserIP(r); ip != "127.0.0.1" {
		log.Printf("geoIP access from %s rejected\n", ip)
		return nil
	}

	if r.URL.Path == "/:geo-ip-put" {
		ip := r.URL.Query().Get("ip")
		value := r.URL.Query().Get("value")
		if !validIP4(ip) || value == "" {
			log.Printf("invalid geoIP put query: ip='%s' value='%s'\n", ip, value)
			return respERR
		}

		location, err := NewLocationFromString(ip, value)
		if err != nil {
			return respERR
		}

		err = c.Put(ip, location)
		if err != nil {
			return respERR
		}
		return respOK
	}
	if r.URL.Path == "/:geo-ip-get" {
		ip := r.URL.Query().Get("ip")
		if !validIP4(ip) {
			return respERR
		}

		result, err := c.Read(ip)
		if result == nil || err != nil {
			return respERR
		}
		return &routing.Cadre{Body: []byte(result.String())}
	}
	return nil
}

func validIP4(ipAddress string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(strings.Trim(ipAddress, " "))
}
