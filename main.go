package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chubin/wttr.in/internal/cache"
	"github.com/chubin/wttr.in/internal/config"
	"github.com/chubin/wttr.in/internal/generate"
	"github.com/chubin/wttr.in/internal/ip"
	"github.com/chubin/wttr.in/internal/location"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/renderer"
	"github.com/chubin/wttr.in/internal/renderer/oneline"
	"github.com/chubin/wttr.in/internal/server"
	"github.com/chubin/wttr.in/internal/spec"
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
	// Configuring Renderers.
	////////////////////////////

	rendererMap := map[string]weather.Renderer{
		"v1":   &renderer.V1Renderer{},
		"v2":   &renderer.V2Renderer{},
		"j1":   &renderer.J1Renderer{},
		"j2":   &renderer.J2Renderer{},
		"line": oneline.NewOnelineRenderer(),
	}

	////////////////////////////
	// Configuring IP Locators.
	////////////////////////////
	ipLocators := []weather.IPLocator{}
	if cfg.IP.GeoIP2 != "" {
		geoIP2, err := ip.NewIPLocatorGeoIP2(cfg.IP.GeoIP2)
		if err != nil {
			log.Fatalln("geoip2 initalization error:", err)
		}
		ipLocators = append(ipLocators, geoIP2)
	}

	ipCache, err := ip.NewCache(cfg.IP)
	if err != nil {
		log.Fatalln(err)
	}
	ipLocators = append(ipLocators, ip.NewIPCacheLocator(ipCache))

	spec, err := spec.LoadSpecFromAssets()
	if err != nil {
		log.Fatalln("error loading wttr.in options description: ", err)
	}

	lruCache, err := cache.NewLRU(cfg.Cache)
	if err != nil {
		log.Fatalln("error creating lru cache: ", err)
	}

	requestLogger := logging.NewRequestLogger(
		cfg.Logging.AccessLog,
		time.Duration(cfg.Logging.Interval)*time.Second,
	)

	ws := weather.NewWeatherService(
		weather.NewWeatherClient(cfg.Weather.WWO),
		weather.NewCacheLocator(locationCache),
		ipLocators,
		weather.NewQueryParser(spec),
		lruCache,
		requestLogger,
		uplink.NewUplinkProcessor(cfg.Uplink),
		rendererMap,
	)

	return server.Serve(&cfg.Server, &cfg.Logging, ws)
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
	// Optional: Set output format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Check the remaining arguments after flag parsing
	if len(flag.Args()) < 1 {
		logrus.Error("Usage: CMD {gen|srv CONFIG}")
		os.Exit(1)
	}

	// Use flag.Args() instead of os.Args to access non-flag arguments
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
