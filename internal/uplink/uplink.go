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
		uplinkResponse, err = getUplink(r, client, location, opts)
	}

	return uplinkRoute, uplinkResponse, err
}

// getUplink forwards the incoming request to one of the backend uplink servers
// and returns a cacheable response.
//
// It enriches the request with location and options data so the backend
// can render the weather without needing to re-parse everything.
func getUplink(req *http.Request, client *http.Client, location *domain.Location, opts *options.Options) (*domain.CacheEntry, error) {
	// Build the target URL using the original Host and RequestURI.
	// This preserves query parameters, path, etc.
	targetURL := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)

	// Create a new request that will be sent to the uplink backend.
	proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, targetURL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy all original headers.
	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	// Ensure we forward the real client IP.
	if proxyReq.Header.Get("X-Forwarded-For") == "" {
		proxyReq.Header.Set("X-Forwarded-For", ipFromAddr(req.RemoteAddr))
	}

	// Attach location data (geolocation info) as JSON header.
	if location != nil {
		locationJSON, err := json.Marshal(location)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal location: %w", err)
		}
		proxyReq.Header.Set("X-Location", string(locationJSON))
	}

	// Attach parsed options as JSON header.
	if opts != nil {
		optionsJSON, err := json.Marshal(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal options: %w", err)
		}
		proxyReq.Header.Set("X-Options", string(optionsJSON))
	}

	// Perform the request to the backend.
	res, err := client.Do(proxyReq)
	if err != nil {
		return nil, fmt.Errorf("uplink request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read uplink response body: %w", err)
	}

	// Return a CacheEntry with randomized expiration (17–25 minutes) to
	// spread out cache invalidation and reduce thundering herd on backends.
	return &domain.CacheEntry{
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     res.Header.Clone(), // clone to avoid sharing mutable header map
		StatusCode: res.StatusCode,
	}, nil
}

// checkURLForPNG returns true if the request is for a .png image
// but not for files in the /files/ directory.
func checkURLForPNG(r *http.Request) bool {
	urlStr := r.URL.String()
	return strings.Contains(urlStr, ".png") && !strings.Contains(urlStr, "/files/")
}

// randInt returns a random integer in the range [min, max).
// Note: rand.Seed should be called once at program startup for better randomness.
func randInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min)
}

// ipFromAddr extracts the IP address from an "IP:PORT" string.
// Returns the original string if no colon is found.
func ipFromAddr(s string) string {
	if pos := strings.LastIndex(s, ":"); pos != -1 {
		return s[:pos]
	}
	return s
}
