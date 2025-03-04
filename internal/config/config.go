package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/chubin/wttr.in/internal/types"
	"github.com/chubin/wttr.in/internal/util"
)

// Config of the program.
type Config struct {
	Cache
	Geo
	Logging
	Server
	Uplink
}

// Logging configuration.
type Logging struct {
	// AccessLog path.
	AccessLog string `yaml:"accessLog,omitempty"`

	// ErrorsLog path.
	ErrorsLog string `yaml:"errorsLog,omitempty"`

	// Interval between access log flushes, in seconds.
	Interval int `yaml:"interval,omitempty"`
}

// Server configuration.
type Server struct {
	// PortHTTP is port where HTTP server must listen.
	// If 0, HTTP is disabled.
	PortHTTP int `yaml:"portHttp,omitempty"`

	// PortHTTP is port where the HTTPS server must listen.
	// If 0, HTTPS is disabled.
	PortHTTPS int `yaml:"portHttps,omitempty"`

	// TLSCertFile contains path to cert file for TLS Server.
	TLSCertFile string `yaml:"tlsCertFile,omitempty"`

	// TLSCertFile contains path to key file for TLS Server.
	TLSKeyFile string `yaml:"tlsKeyFile,omitempty"`
}

// Uplink configuration.
type Uplink struct {
	// Address1 contains address of the uplink server in form IP:PORT
	// for format=j1 queries.
	Address1 string `yaml:"address1,omitempty"`

	// Address2 contains address of the uplink server in form IP:PORT
	// for format=* queries.
	Address2 string `yaml:"address2,omitempty"`

	// Address3 contains address of the uplink server in form IP:PORT
	// for all other queries.
	Address3 string `yaml:"address3,omitempty"`

	// Address4 contains address of the uplink server in form IP:PORT
	// for PNG queries.
	Address4 string `yaml:"address4,omitempty"`

	// Timeout for upstream queries.
	Timeout int `yaml:"timeout,omitempty"`

	// PrefetchInterval contains time (in milliseconds) indicating,
	// how long the prefetch procedure should take.
	PrefetchInterval int `yaml:"prefetchInterval,omitempty"`
}

// Cache configuration.
type Cache struct {
	// Size of the main cache.
	Size int `yaml:"size,omitempty"`
}

// Geo contains geolocation configuration.
type Geo struct {
	// IPCache contains the path to the IP Geodata cache.
	IPCache string `yaml:"ipCache,omitempty"`

	// IPCacheDB contains the path to the SQLite DB with the IP Geodata cache.
	IPCacheDB string `yaml:"ipCacheDb,omitempty"`

	IPCacheType types.CacheType `yaml:"ipCacheType,omitempty"`

	// LocationCache contains the path to the Location Geodata cache.
	LocationCache string `yaml:"locationCache,omitempty"`

	// LocationCacheDB contains the path to the SQLite DB with the Location Geodata cache.
	LocationCacheDB string `yaml:"locationCacheDb,omitempty"`

	LocationCacheType types.CacheType `yaml:"locationCacheType,omitempty"`

	Nominatim []Nominatim `yaml:"nominatim"`
}

type Nominatim struct {
	Name string

	// Type describes the type of the location service.
	// Supported types: iq.
	Type string

	URL string

	Token string
}

// Default contains the default configuration.
func Default() *Config {
	return &Config{
		Cache{
			Size: 12800,
		},
		Geo{
			IPCache:           "/wttr.in/cache/ip2l",
			IPCacheDB:         "/wttr.in/cache/geoip.db",
			IPCacheType:       types.CacheTypeDB,
			LocationCache:     "/wttr.in/cache/loc",
			LocationCacheDB:   "/wttr.in/cache/geoloc.db",
			LocationCacheType: types.CacheTypeDB,
			Nominatim: []Nominatim{
				{
					Name:  "locationiq",
					Type:  "iq",
					URL:   "https://eu1.locationiq.com/v1/search",
					Token: os.Getenv("NOMINATIM_LOCATIONIQ"),
				},
				{
					Name:  "opencage",
					Type:  "opencage",
					URL:   "https://api.opencagedata.com/geocode/v1/json",
					Token: os.Getenv("NOMINATIM_OPENCAGE"),
				},
			},
		},
		Logging{
			AccessLog: "/wttr.in/log/access.log",
			ErrorsLog: "/wttr.in/log/errors.log",
			Interval:  300,
		},
		Server{
			PortHTTP:    8083,
			PortHTTPS:   8084,
			TLSCertFile: "/wttr.in/etc/fullchain.pem",
			TLSKeyFile:  "/wttr.in/etc/privkey.pem",
		},
		Uplink{
			Address1:         "127.0.0.1:9002",
			Address2:         "127.0.0.1:9002",
			Address3:         "127.0.0.1:9002",
			Address4:         "127.0.0.1:9002",
			Timeout:          30,
			PrefetchInterval: 300,
		},
	}
}

// Load config from file.
func Load(filename string) (*Config, error) {
	var (
		config Config
		data   []byte
		err    error
	)

	data, err = os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = util.YamlUnmarshalStrict(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Dump() []byte {
	data, err := yaml.Marshal(c)
	if err != nil {
		// should never happen.
		log.Fatalln("config.Dump():", err)
	}

	return data
}
