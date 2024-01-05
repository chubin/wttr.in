package processor

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"github.com/chubin/wttr.in/internal/config"
	geoip "github.com/chubin/wttr.in/internal/geo/ip"
	geoloc "github.com/chubin/wttr.in/internal/geo/location"
	"github.com/chubin/wttr.in/internal/routing"
	"github.com/chubin/wttr.in/internal/stats"
	"github.com/chubin/wttr.in/internal/util"
)

// plainTextAgents contains signatures of the plain-text agents.
func plainTextAgents() []string {
	return []string{
		"curl",
		"httpie",
		"lwp-request",
		"wget",
		"python-httpx",
		"python-requests",
		"openbsd ftp",
		"powershell",
		"fetch",
		"aiohttp",
		"http_get",
		"xh",
		"nushell",
	}
}

type ResponseWithHeader struct {
	InProgress bool      // true if the request is being processed
	Expires    time.Time // expiration time of the cache entry

	Body       []byte
	Header     http.Header
	StatusCode int // e.g. 200
}

// RequestProcessor handles incoming requests.
type RequestProcessor struct {
	peakRequest30     sync.Map
	peakRequest60     sync.Map
	lruCache          *lru.Cache
	stats             *stats.Stats
	router            routing.Router
	upstreamTransport *http.Transport
	config            *config.Config
	geoIPCache        *geoip.Cache
	geoLocation       *geoloc.Cache
}

// NewRequestProcessor returns new RequestProcessor.
func NewRequestProcessor(config *config.Config) (*RequestProcessor, error) {
	lruCache, err := lru.New(config.Cache.Size)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{
		Timeout:   time.Duration(config.Uplink.Timeout) * time.Second,
		KeepAlive: time.Duration(config.Uplink.Timeout) * time.Second,
		DualStack: true,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, config.Uplink.Address)
		},
	}

	geoCache, err := geoip.NewCache(config)
	if err != nil {
		return nil, err
	}

	geoLocation, err := geoloc.NewCache(config)
	if err != nil {
		return nil, err
	}

	rp := &RequestProcessor{
		lruCache:          lruCache,
		stats:             stats.New(),
		upstreamTransport: transport,
		config:            config,
		geoIPCache:        geoCache,
		geoLocation:       geoLocation,
	}

	// Initialize routes.
	rp.router.AddPath("/:stats", rp.stats)
	rp.router.AddPath("/:geo-ip-get", rp.geoIPCache)
	rp.router.AddPath("/:geo-ip-put", rp.geoIPCache)
	rp.router.AddPath("/:geo-location", rp.geoLocation)

	return rp, nil
}

// Start starts async request processor jobs, such as peak handling.
func (rp *RequestProcessor) Start() error {
	return rp.startPeakHandling()
}

func (rp *RequestProcessor) ProcessRequest(r *http.Request) (*ResponseWithHeader, error) {
	var (
		response *ResponseWithHeader
		ip       = util.ReadUserIP(r)
	)

	if ip != "127.0.0.1" {
		rp.stats.Inc("total")
	}

	// Main routing logic.
	if rh := rp.router.Route(r); rh != nil {
		result := rh.Response(r)
		if result != nil {
			return fromCadre(result), nil
		}
	}

	if resp, ok := redirectInsecure(r); ok {
		rp.stats.Inc("redirects")

		return resp, nil
	}

	if dontCache(r) {
		rp.stats.Inc("uncached")

		return get(r, rp.upstreamTransport)
	}

	// processing cached request
	cacheDigest := getCacheDigest(r)

	rp.savePeakRequest(cacheDigest, r)

	response = rp.processRequestFromCache(r)
	if response != nil {
		return response, nil
	}

	return rp.processUncachedRequest(r)
}

// processRequestFromCache processes requests using the cache.
// If no entry in cache found, nil is returned.
func (rp *RequestProcessor) processRequestFromCache(r *http.Request) *ResponseWithHeader {
	var (
		cacheEntry  ResponseWithHeader
		cacheDigest = getCacheDigest(r)
		ok          bool
	)

	cacheBody, _ := rp.lruCache.Get(cacheDigest)
	cacheEntry, ok = cacheBody.(ResponseWithHeader)
	if !ok {
		return nil
	}

	// if after all attempts we still have no answer,
	// we try to make the query on our own
	for attempts := 0; attempts < 300; attempts++ {
		if !ok || !cacheEntry.InProgress {
			break
		}
		time.Sleep(30 * time.Millisecond)
		cacheBody, _ = rp.lruCache.Get(cacheDigest)
		v, ok := cacheBody.(ResponseWithHeader)
		if ok {
			cacheEntry = v
		}
	}
	if cacheEntry.InProgress {
		log.Printf("TIMEOUT: %s\n", cacheDigest)
	}
	if ok && !cacheEntry.InProgress && cacheEntry.Expires.After(time.Now()) {
		rp.stats.Inc("cache1")

		return &cacheEntry
	}

	return nil
}

// processUncachedRequest processes requests that were not found in the cache.
func (rp *RequestProcessor) processUncachedRequest(r *http.Request) (*ResponseWithHeader, error) {
	var (
		cacheDigest = getCacheDigest(r)
		ip          = util.ReadUserIP(r)
		response    *ResponseWithHeader
		err         error
	)

	// Response was not found in cache.
	// Starting real handling.
	format := r.URL.Query().Get("format")
	if len(format) != 0 {
		rp.stats.Inc("format")
		if format == "j1" {
			rp.stats.Inc("format=j1")
		}
	}

	// Count, how many IP addresses are known.
	_, err = rp.geoIPCache.Read(ip)
	if err == nil {
		rp.stats.Inc("geoip")
	}

	// Indicate, that the request is being handled.
	rp.lruCache.Add(cacheDigest, ResponseWithHeader{InProgress: true})

	response, err = get(r, rp.upstreamTransport)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 || response.StatusCode == 304 || response.StatusCode == 404 {
		rp.lruCache.Add(cacheDigest, *response)
	} else {
		log.Printf("REMOVE: %d response for %s from cache\n", response.StatusCode, cacheDigest)
		rp.lruCache.Remove(cacheDigest)
	}

	return response, nil
}

// getCacheDigest is an implementation of the cache.get_signature of original wttr.in.
func getCacheDigest(req *http.Request) string {
	userAgent := req.Header.Get("User-Agent")

	queryHost := req.Host
	queryString := req.RequestURI

	clientIPAddress := util.ReadUserIP(req)

	lang := req.Header.Get("Accept-Language")

	return fmt.Sprintf("%s:%s%s:%s:%s", userAgent, queryHost, queryString, clientIPAddress, lang)
}

// dontCache returns true if req should not be cached.
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
// `proxy_set_header   X-Forwarded-Proto $scheme;`.
func redirectInsecure(req *http.Request) (*ResponseWithHeader, bool) {
	if isPlainTextAgent(req.Header.Get("User-Agent")) {
		return nil, false
	}

	if req.TLS != nil || strings.ToLower(req.Header.Get("X-Forwarded-Proto")) == "https" {
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

	return &ResponseWithHeader{
		InProgress: false,
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     http.Header{"Location": []string{target}},
		StatusCode: 301,
	}, true
}

// isPlainTextAgent returns true if userAgent is a plain-text agent.
func isPlainTextAgent(userAgent string) bool {
	userAgentLower := strings.ToLower(userAgent)
	for _, signature := range plainTextAgents() {
		if strings.Contains(userAgentLower, signature) {
			return true
		}
	}

	return false
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// ipFromAddr returns IP address from a ADDR:PORT pair.
func ipFromAddr(s string) string {
	pos := strings.LastIndex(s, ":")
	if pos == -1 {
		return s
	}

	return s[:pos]
}

// fromCadre converts Cadre into a responseWithHeader.
func fromCadre(cadre *routing.Cadre) *ResponseWithHeader {
	return &ResponseWithHeader{
		Body:       cadre.Body,
		Expires:    cadre.Expires,
		StatusCode: 200,
		InProgress: false,
	}
}
