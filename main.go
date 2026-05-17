package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/chubin/wttr.in/internal/cache"
	"github.com/chubin/wttr.in/internal/config"
	"github.com/chubin/wttr.in/internal/defs"
	"github.com/chubin/wttr.in/internal/formatter"
	"github.com/chubin/wttr.in/internal/generate"
	"github.com/chubin/wttr.in/internal/ip"
	"github.com/chubin/wttr.in/internal/location"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/query"
	"github.com/chubin/wttr.in/internal/renderer"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
	"github.com/chubin/wttr.in/internal/renderer/subprocess"
	v1 "github.com/chubin/wttr.in/internal/renderer/v1"
	v2 "github.com/chubin/wttr.in/internal/renderer/v2"
	"github.com/chubin/wttr.in/internal/server"
	"github.com/chubin/wttr.in/internal/translate"
	"github.com/chubin/wttr.in/internal/uplink"
	"github.com/chubin/wttr.in/internal/weather"
)

var debug = true

func srv(configFile string) error {
	cfg, err := config.LoadFromYAML(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if debug {
		configYAML, err := cfg.MarshalYAML()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("loaded config:\n" + string(configYAML))
	}

	locationCache, err := location.NewCache(cfg.Geo)
	if err != nil {
		log.Fatalln(err)
	}

	////////////////////////////
	// Configuring Renderers
	////////////////////////////
	rendererMap := map[string]weather.Renderer{
		"v1":         v1.NewV1Renderer(),
		"v2":         v2.NewV2Renderer(),
		"v2d":        v2.NewV2Renderer(),
		"v2n":        v2.NewV2Renderer(),
		"p1":         renderer.NewPrometheusRenderer(),
		"j1":         &renderer.J1Renderer{},
		"j2":         &renderer.J2Renderer{},
		"line":       oneline.NewOnelineRenderer(),
		"page":       renderer.NewPageRenderer(),
		"subprocess": subprocess.NewRenderer(cfg.Renderer.Subprocess),
	}

	////////////////////////////
	// Configuring Formatters
	////////////////////////////
	htmlFormatter, err := formatter.NewHTMLFormatter()
	if err != nil {
		return fmt.Errorf("html formatter creation error: %w", err)
	}

	formatterMap := map[string]weather.Formatter{
		"text": &formatter.TextFormatter{},
		"html": htmlFormatter,
		"png":  formatter.NewPNGFormatter(),
	}

	////////////////////////////
	// Configuring IP Locators
	////////////////////////////
	ipLocators := []weather.IPLocator{}
	if cfg.IP.GeoIP2 != "" {
		geoIP2, err := ip.NewIPLocatorGeoIP2(cfg.IP.GeoIP2)
		if err != nil {
			log.Fatalln("geoip2 initialization error:", err)
		}
		ipLocators = append(ipLocators, geoIP2)
	}

	ipCache, err := ip.NewCache(cfg.IP)
	if err != nil {
		log.Fatalln(err)
	}
	ipLocators = append(ipLocators, ip.NewIPCacheLocator(ipCache))

	defs, err := defs.LoadDefsFromAssets()
	if err != nil {
		log.Fatalln("error loading wttr.in options description: ", err)
	}

	requestLogger := logging.NewRequestLogger(
		cfg.Logging.AccessLog,
		time.Duration(cfg.Logging.Interval)*time.Second,
	)

	localizer := translate.NewBundle(assets.FS)

	// ==================== MULTI-LAYER CACHE SETUP ====================

	// 1. Responses Cache - final rendered output (usually LRU / in-memory)
	var responsesCacher weather.Cacher
	if cfg.Cache.Responses.IsEnabled() {
		if cfg.Cache.Responses.IsLRU() || cfg.Cache.Responses.Type == "" {
			responsesCacher, err = cache.NewLRU(cfg.Cache.Responses)
			if err != nil {
				log.Fatalln("failed to create responses LRU cache:", err)
			}
		} else {
			// fallback
			responsesCacher, err = cache.NewLRU(cfg.Cache.Responses)
			if err != nil {
				log.Fatalln("failed to create responses cache:", err)
			}
		}
	} else {
		responsesCacher = cache.NewNoOpCacher()
	}

	// 2. Weather Cache - raw data from upstream (WWO, etc.)
	var weatherCacher weather.Cacher
	if cfg.Cache.Weather.IsEnabled() {
		if cfg.Cache.Weather.IsDisk() {
			weatherCacher, err = cache.NewDiskCacher(cfg.Cache.Weather.Dir)
			if err != nil {
				log.Fatalln("failed to create weather disk cache:", err)
			}
		} else if cfg.Cache.Weather.IsLRU() {
			weatherCacher, err = cache.NewLRU(cfg.Cache.Weather)
			if err != nil {
				log.Fatalln("failed to create weather LRU cache:", err)
			}
		} else {
			weatherCacher, err = cache.NewLRU(cfg.Cache.Weather) // fallback
			if err != nil {
				log.Fatalln("failed to create weather cache:", err)
			}
		}
	} else {
		weatherCacher = cache.NewNoOpCacher()
	}

	// 3. Create cached weather client
	rawClient := weather.NewWeatherClient(cfg.Weather.WWO)

	cachedWeatherClient := weather.NewCachedWeatherClient(
		rawClient,
		weatherCacher,
		cfg.Cache.Weather.TTL,
	)

	// ==================== END CACHE SETUP ====================

	ws := weather.NewWeatherService(
		cachedWeatherClient,
		weather.NewCacheLocator(locationCache),
		ipLocators,
		query.NewQueryParser(defs),
		responsesCacher,
		requestLogger,
		uplink.NewUplinkProcessor(cfg.Uplink),
		rendererMap,
		formatterMap,
		localizer,
	)

	return server.Serve(cfg.Server, cfg.Logging, ws)
}

func main() {
	flagDebug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	debug = *flagDebug

	// Configure logrus
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if len(flag.Args()) < 1 {
		logrus.Error("Usage: CMD {gen|srv CONFIG}")
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "srv":
		var configFile string
		if len(flag.Args()) > 1 {
			configFile = flag.Args()[1]
		} else {
			log.Fatalln("usage: CMD srv CONFIG")
		}

		logrus.Info("Starting server...")
		err := srv(configFile)
		if err != nil {
			log.Fatalln(err)
		}

	case "gen":
		logrus.Info("Generating options and parser...")
		err := generate.GenerateOptionsAndParser()
		if err != nil {
			logrus.Error(err)
		}

	default:
		logrus.Error("Invalid command. Usage: CMD {gen|srv}")
		os.Exit(1)
	}
}
