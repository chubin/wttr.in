package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

func processRequest(r *http.Request) responseWithHeader {
	var response responseWithHeader

	if response, ok := redirectInsecure(r); ok {
		return *response
	}

	if dontCache(r) {
		return get(r)
	}

	cacheDigest := getCacheDigest(r)

	foundInCache := false

	savePeakRequest(cacheDigest, r)

	cacheBody, ok := lruCache.Get(cacheDigest)
	if ok {
		cacheEntry := cacheBody.(responseWithHeader)

		// if after all attempts we still have no answer,
		// we try to make the query on our own
		for attempts := 0; attempts < 300; attempts++ {
			if !ok || !cacheEntry.InProgress {
				break
			}
			time.Sleep(30 * time.Millisecond)
			cacheBody, ok = lruCache.Get(cacheDigest)
			cacheEntry = cacheBody.(responseWithHeader)
		}
		if cacheEntry.InProgress {
			log.Printf("TIMEOUT: %s\n", cacheDigest)
		}
		if ok && !cacheEntry.InProgress && cacheEntry.Expires.After(time.Now()) {
			response = cacheEntry
			foundInCache = true
		}
	}

	if !foundInCache {
		lruCache.Add(cacheDigest, responseWithHeader{InProgress: true})
		response = get(r)
		if response.StatusCode == 200 || response.StatusCode == 304 || response.StatusCode == 404 {
			lruCache.Add(cacheDigest, response)
		} else {
			log.Printf("REMOVE: %d response for %s from cache\n", response.StatusCode, cacheDigest)
			lruCache.Remove(cacheDigest)
		}
	}
	return response
}

func get(req *http.Request) responseWithHeader {

	client := &http.Client{}

	queryURL := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, queryURL, req.Body)
	if err != nil {
		log.Printf("Request: %s\n", err)
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
	if err != nil {
		log.Println(err)
	}

	return responseWithHeader{
		InProgress: false,
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}
}

// implementation of the cache.get_signature of original wttr.in
func getCacheDigest(req *http.Request) string {

	userAgent := req.Header.Get("User-Agent")

	queryHost := req.Host
	queryString := req.RequestURI

	clientIPAddress := readUserIP(req)

	lang := req.Header.Get("Accept-Language")

	return fmt.Sprintf("%s:%s%s:%s:%s", userAgent, queryHost, queryString, clientIPAddress, lang)
}

// return true if request should not be cached
func dontCache(req *http.Request) bool {

	// dont cache cyclic requests
	loc := strings.Split(req.RequestURI, "?")[0]
	return strings.Contains(loc, ":")
}

// redirectInsecure returns redirection response, and bool value, if redirection was needed,
// if the query comes from a browser, and it is insecure.
//
// Insecure queries are marked by the frontend web server
// with X-Forwarded-Proto header:
//
//    proxy_set_header   X-Forwarded-Proto $scheme;
//
//
func redirectInsecure(req *http.Request) (*responseWithHeader, bool) {
	if isPlainTextAgent(req.Header.Get("User-Agent")) {
		return nil, false
	}

	if strings.ToLower(req.Header.Get("X-Forwarded-Proto")) == "https" {
		return nil, false
	}

	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	body := []byte(fmt.Sprintf(`<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>301 Moved</TITLE></HEAD><BODY>
<H1>301 Moved</H1>
The document has moved
<A HREF="%s">here</A>.
</BODY></HTML>
`, target))

	return &responseWithHeader{
		InProgress: false,
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     http.Header{"Location": []string{target}},
		StatusCode: 301,
	}, true
}

// isPlainTextAgent returns true if userAgent is a plain-text agent
func isPlainTextAgent(userAgent string) bool {
	userAgentLower := strings.ToLower(userAgent)
	for _, signature := range plainTextAgents {
		if strings.Contains(userAgentLower, signature) {
			return true
		}
	}
	return false
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
			log.Printf("ERROR: userip: %q is not IP:port\n", IPAddress)
		}
	}
	return IPAddress
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
