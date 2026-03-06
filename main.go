package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chubin/wttr.go/internal/generate"
	"github.com/chubin/wttr.go/internal/location"
	"github.com/chubin/wttr.go/internal/options"
	"github.com/chubin/wttr.go/internal/weather"
)

func srv() {
	locationCache, err := location.NewCache(
		&location.Config{
			LocationCacheType: "",
			LocationCacheDB:   "",
			LocationCache:     "",
			IPCacheType:       "",
			NominatimServers: []struct {
				Name  string
				Type  string
				URL   string
				Token string
			}{
				{
					"",
					"",
					"",
					"",
				},
			},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	wttrInOptions, err := options.NewFromFile("spec/options/options.yaml")
	if err != nil {
		log.Fatalln("error loading wttr.in options description: ", err)
	}

	ws := weather.NewWeatherService(
		weather.NewWeatherClient(fmt.Sprintf(
			"http://127.0.0.1:5001/premium/v1/weather.ashx?key=%s&q={lat},{long}&format=json&num_of_days=3&includelocation=yes",
			os.Getenv("PROXY_KEY"),
		)),
		weather.NewCacheLocator(locationCache),
		weather.NewIPCacheLocator(nil),
		weather.NewQueryParser(wttrInOptions),
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: CMD {gen|srv}")
		os.Exit(0)
	}

	if os.Args[1] == "srv" {
		srv()
	}
	if os.Args[1] == "gen" {
		err := generate.GenerateOptions()
		if err != nil {
			fmt.Println(err)
		}
	}
}
