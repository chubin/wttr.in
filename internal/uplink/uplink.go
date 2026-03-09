package uplink

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
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
	uplinkTransport1 *http.Transport
	uplinkTransport2 *http.Transport
	uplinkTransport3 *http.Transport
	uplinkTransport4 *http.Transport
}

func NewUplinkProcessor(cfg Config) *UplinkProcessor {
	return &UplinkProcessor{}
}

func (s *UplinkProcessor) Route(
	opts *query.Options, r *http.Request, ipData *weather.IPData, location *weather.Location,
) (bool, *weather.CacheEntry, error) {
	return false, nil, nil
}

func getAny(req *http.Request, tr1, tr2, tr3, tr4 *http.Transport) (*weather.CacheEntry, error) {
	uri := strings.ReplaceAll(req.URL.RequestURI(), "%", "%25")

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	format := u.Query().Get("format")

	if format == "j1" {
		return getJ1(req, tr1)
	} else if format != "" {
		return getFormat(req, tr2)
	}

	if checkURLForPNG(req) {
		return getDefault(req, tr4)
	}

	return getDefault(req, tr3)
}

func getJ1(req *http.Request, transport *http.Transport) (*weather.CacheEntry, error) {
	return getUplink(req, transport)
}

func getFormat(req *http.Request, transport *http.Transport) (*weather.CacheEntry, error) {
	return getUplink(req, transport)
}

func getDefault(req *http.Request, transport *http.Transport) (*weather.CacheEntry, error) {
	return getUplink(req, transport)
}

func getUplink(req *http.Request, transport *http.Transport) (*weather.CacheEntry, error) {
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
