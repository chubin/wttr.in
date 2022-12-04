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
	"github.com/chubin/wttr.in/internal/routing"
	"github.com/chubin/wttr.in/internal/stats"
	"github.com/chubin/wttr.in/internal/util"
)

// plainTextAgents contains signatures of the plain-text agents.
var plainTextAgents = []string{
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
}

type responseWithHeader struct {
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

	rp := &RequestProcessor{
		lruCache:          lruCache,
		stats:             stats.New(),
		upstreamTransport: transport,
		config:            config,
		geoIPCache:        geoip.NewCache(config),
	}

	// Initialize routes.
	rp.router.AddPath("/:stats", rp.stats)

	return rp, nil
}

// Start starts async request processor jobs, such as peak handling.
func (rp *RequestProcessor) Start() {
	rp.startPeakHandling()
}

func (rp *RequestProcessor) ProcessRequest(r *http.Request) (*responseWithHeader, error) {
	var (
		response *responseWithHeader
		err      error
	)

	rp.stats.Inc("total")

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

	cacheDigest := getCacheDigest(r)

	foundInCache := false

	rp.savePeakRequest(cacheDigest, r)

	cacheBody, ok := rp.lruCache.Get(cacheDigest)
	if ok {
		rp.stats.Inc("cache1")
		cacheEntry := cacheBody.(responseWithHeader)

		// if after all attempts we still have no answer,
		// we try to make the query on our own
		for attempts := 0; attempts < 300; attempts++ {
			if !ok || !cacheEntry.InProgress {
				break
			}
			time.Sleep(30 * time.Millisecond)
			cacheBody, ok = rp.lruCache.Get(cacheDigest)
			if ok && cacheBody != nil {
				cacheEntry = cacheBody.(responseWithHeader)
			}
		}
		if cacheEntry.InProgress {
			log.Printf("TIMEOUT: %s\n", cacheDigest)
		}
		if ok && !cacheEntry.InProgress && cacheEntry.Expires.After(time.Now()) {
			response = &cacheEntry
			foundInCache = true
		}
	}

	if !foundInCache {
		// Handling query.
		format := r.URL.Query().Get("format")
		if len(format) != 0 {
			rp.stats.Inc("format")
			if format == "j1" {
				rp.stats.Inc("format=j1")
			}
		}

		// How many IP addresses are known.
		ip := util.ReadUserIP(r)
		_, err = rp.geoIPCache.Read(ip)
		if err == nil {
			rp.stats.Inc("geoip")
		}

		rp.lruCache.Add(cacheDigest, responseWithHeader{InProgress: true})

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
	}
	return response, nil
}

func get(req *http.Request, transport *http.Transport) (*responseWithHeader, error) {

	client := &http.Client{
		Transport: transport,
	}

	queryURL := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, queryURL, req.Body)
	if err != nil {
		return nil, err
	}

	// proxyReq.Header.Set("Host", req.Host)
	// proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	if proxyReq.Header.Get("X-Forwarded-For") == "" {
		proxyReq.Header.Set("X-Forwarded-For", ipFromAddr(req.RemoteAddr))
	}

	res, err := client.Do(proxyReq)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &responseWithHeader{
		InProgress: false,
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}, nil
}

// implementation of the cache.get_signature of original wttr.in
func getCacheDigest(req *http.Request) string {

	userAgent := req.Header.Get("User-Agent")

	queryHost := req.Host
	queryString := req.RequestURI

	clientIPAddress := util.ReadUserIP(req)

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
func fromCadre(cadre *routing.Cadre) *responseWithHeader {
	return &responseWithHeader{
		Body:       cadre.Body,
		Expires:    cadre.Expires,
		StatusCode: 200,
		InProgress: false,
	}

}
