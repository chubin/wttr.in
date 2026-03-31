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

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/options"
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
	client1 *http.Client
	client2 *http.Client
	client3 *http.Client
	client4 *http.Client
}

func NewUplinkProcessor(cfg Config) *UplinkProcessor {
	dialer := &net.Dialer{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		KeepAlive: time.Duration(cfg.Timeout) * time.Second,
		DualStack: true,
	}

	mkClient := func(addr string) *http.Client {
		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
					return dialer.DialContext(ctx, network, addr)
				},
				MaxIdleConnsPerHost: 32, // tune this!
				IdleConnTimeout:     90 * time.Second,
			},
			Timeout: 15 * time.Second, // safety net
		}
	}

	return &UplinkProcessor{
		client1: mkClient(cfg.Address1),
		client2: mkClient(cfg.Address2),
		client3: mkClient(cfg.Address3),
		client4: mkClient(cfg.Address4),
	}
}

func (p *UplinkProcessor) Route(
	opts *options.Options, r *http.Request, ipData *domain.IPData, location *domain.Location,
) (bool, *domain.CacheEntry, error) {
	var (
		uplinkRoute    bool = true
		uplinkResponse *domain.CacheEntry
		err            error
		client         *http.Client
	)

	//////////////////////////////////////////

	// Views that are not processed by the uplink.
	if opts.View == "line" || opts.View == "j1" || opts.View == "j2" {
		return false, nil, nil
	}

	if checkURLForPNG(r) {
		client = p.client4
	} else if opts.View == "v1" || opts.View == "files" {
		client = p.client3
	} else if opts.View == "v2" || opts.View == "p1" {
		client = p.client2
	} else {
		// The rest goes to the client1 (should be empty).
		client = p.client1
	}
	//////////////////////////////////////////

	if uplinkRoute {
		uplinkResponse, err = getUplink(r, client, location)
	}

	return uplinkRoute, uplinkResponse, err
}

func getUplink(req *http.Request, client *http.Client, location *domain.Location) (*domain.CacheEntry, error) {
	// client := &http.Client{
	// 	Transport: transport,
	// }

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

	return &domain.CacheEntry{
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
