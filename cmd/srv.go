package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru"
)

var lruCache *lru.Cache

type ResponseWithHeader struct {
	Body       []byte
	Header     http.Header
	StatusCode int // e.g. 200

}

func init() {
	var err error
	lruCache, err = lru.New(12800)
	if err != nil {
		panic(err)
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		addr = "127.0.0.1:8002"
		return dialer.DialContext(ctx, network, addr)
	}

}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
		var err error
		IPAddress, _, err = net.SplitHostPort(IPAddress)
		if err != nil {
			fmt.Printf("userip: %q is not IP:port\n", IPAddress)
		}
	}
	return IPAddress
}

// implementation of the cache.get_signature of original wttr.in
func findCacheDigest(req *http.Request) string {

	userAgent := req.Header.Get("User-Agent")

	queryHost := req.Host
	queryString := req.RequestURI

	clientIpAddress := readUserIP(req)

	lang := req.Header.Get("Accept-Language")

	now := time.Now()
	secs := now.Unix()
	timestamp := secs / 1000

	return fmt.Sprintf("%s:%s%s:%s:%s:%d", userAgent, queryHost, queryString, clientIpAddress, lang, timestamp)
}

func get(req *http.Request) ResponseWithHeader {

	client := &http.Client{}

	queryURL := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, queryURL, req.Body)
	if err != nil {
		// handle error
	}

	// proxyReq.Header.Set("Host", req.Host)
	// proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	res, err := client.Do(proxyReq)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	return ResponseWithHeader{
		Body:       body,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}
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
		var response ResponseWithHeader

		cacheDigest := findCacheDigest(r)
		cacheBody, ok := lruCache.Get(cacheDigest)
		if ok {
			response = cacheBody.(ResponseWithHeader)
		} else {
			fmt.Println(cacheDigest)
			response = get(r)
			if response.StatusCode == 200 {
				lruCache.Add(cacheDigest, response)
			}
		}
		copyHeader(w.Header(), response.Header)
		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
