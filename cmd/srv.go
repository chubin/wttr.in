package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const uplinkSrvAddr = "127.0.0.1:9002"
const uplinkTimeout = 30
const prefetchInterval = 300
const lruCacheSize = 12800

var lruCache *lru.Cache

type responseWithHeader struct {
	InProgress bool      // true if the request is being processed
	Expires    time.Time // expiration time of the cache entry

	Body       []byte
	Header     http.Header
	StatusCode int // e.g. 200
}

func init() {
	var err error
	lruCache, err = lru.New(lruCacheSize)
	if err != nil {
		panic(err)
	}

	dialer := &net.Dialer{
		Timeout:   uplinkTimeout * time.Second,
		KeepAlive: uplinkTimeout * time.Second,
		DualStack: true,
	}

	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, uplinkSrvAddr)
	}

	initPeakHandling()
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// printStat()
		response := processRequest(r)

		copyHeader(w.Header(), response.Header)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
