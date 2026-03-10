package uplink

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/chubin/wttr.go/internal/query"
	"github.com/chubin/wttr.go/internal/weather"
)

type Config struct {
	Address1         string `yaml:"address1"`
	Address2         string `yaml:"address2"`
	Address3         string `yaml:"address3"`
	Address4         string `yaml:"address4"`
	Timeout          int    `yaml:"timeout"`
	PrefetchInterval int    `yaml:"prefetchInterval"`
}

// UplinkProcessor handles incoming requests.
type UplinkProcessor struct {
	transport1 *http.Transport
	transport2 *http.Transport
	transport3 *http.Transport
	transport4 *http.Transport
}

func NewUplinkProcessor(cfg Config) *UplinkProcessor {
	dialer := &net.Dialer{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		KeepAlive: time.Duration(cfg.Timeout) * time.Second,
		DualStack: true,
	}

	transport1 := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, cfg.Address1)
		},
	}
	transport2 := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, cfg.Address2)
		},
	}
	transport3 := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, cfg.Address3)
		},
	}
	transport4 := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, cfg.Address4)
		},
	}

	return &UplinkProcessor{
		transport1: transport1,
		transport2: transport2,
		transport3: transport3,
		transport4: transport4,
	}
}

func (p *UplinkProcessor) Route(
	opts *query.Options, r *http.Request, ipData *weather.IPData, location *weather.Location,
) (bool, *weather.CacheEntry, error) {
	var (
		uplinkRoute    bool = true
		uplinkResponse *weather.CacheEntry
		err            error
		transport      *http.Transport
	)

	//////////////////////////////////////////

	if checkURLForPNG(r) {
		transport = p.transport4
	} else if opts.View != "" && opts.View != "v1" && opts.View != "j1" && opts.View != "j2" {
		transport = p.transport1
	} else if opts.View == "v1" {
		transport = p.transport2
	} else if opts.View == "v1" {
		transport = p.transport3
	} else {
		uplinkRoute = false
	}
	//////////////////////////////////////////

	if uplinkRoute {
		uplinkResponse, err = getUplink(r, transport, location)
	}

	return uplinkRoute, uplinkResponse, err
}

func getUplink(req *http.Request, transport *http.Transport, location *weather.Location) (*weather.CacheEntry, error) {
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

	locationJson, err := json.Marshal(location)
	if err != nil {
		return nil, err
	}

	proxyReq.Header.Set("X-Location", string(locationJson))

	res, err := client.Do(proxyReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &weather.CacheEntry{
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}, nil
}

func checkURLForPNG(r *http.Request) bool {
	url := r.URL.String()
	return strings.Contains(url, ".png") && !strings.Contains(url, "/files/")
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
