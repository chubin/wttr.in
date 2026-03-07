package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/chubin/wttr.go/internal/cache"
	"github.com/chubin/wttr.go/internal/config"
	"github.com/chubin/wttr.go/internal/generate"
	"github.com/chubin/wttr.go/internal/ip"
	"github.com/chubin/wttr.go/internal/location"
	"github.com/chubin/wttr.go/internal/options"
	"github.com/chubin/wttr.go/internal/weather"
)

var debug = true

func srv() {
	cfg, err := config.LoadFromYAML("config.yaml")
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

	ipCache, err := ip.NewCache(cfg.IP)
	if err != nil {
		log.Fatalln(err)
	}

	wttrInOptions, err := options.NewFromFile("spec/options/options.yaml")
	if err != nil {
		log.Fatalln("error loading wttr.in options description: ", err)
	}

	lruCache, err := cache.NewLRU(cfg.Cache)
	if err != nil {
		log.Fatalln("error creating lru cache: ", err)
	}

	ws := weather.NewWeatherService(
		weather.NewWeatherClient(cfg.Weather.WWO),
		weather.NewCacheLocator(locationCache),
		weather.NewIPCacheLocator(ipCache),
		weather.NewQueryParser(wttrInOptions),
		lruCache,
	)

	// Define routes
	http.HandleFunc("/", ws.WeatherHandler)

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
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
		logrus.Error("Usage: CMD {gen|srv}")
		os.Exit(1)
	}

	// Use flag.Args() instead of os.Args to access non-flag arguments
	switch flag.Args()[0] {
	case "srv":
		logrus.Info("Starting server...")
		srv()
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
