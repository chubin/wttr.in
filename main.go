package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chubin/wttr.go/internal/generate"
	"github.com/chubin/wttr.go/internal/pipeline"
)

func srv() {
	ws := pipeline.NewWeatherService(
		nil,
		pipeline.NewCacheLocator(nil),
		pipeline.NewIPCacheLocator(nil),
		nil,
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
