package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAny(req *http.Request, tr1, tr2, tr3 *http.Transport) (*ResponseWithHeader, error) {
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

	// log.Println(req.URL.Query())
	// log.Println()

	return getDefault(req, tr3)
}

func getJ1(req *http.Request, transport *http.Transport) (*ResponseWithHeader, error) {
	return getUpstream(req, transport)
}

func getFormat(req *http.Request, transport *http.Transport) (*ResponseWithHeader, error) {
	return getUpstream(req, transport)
}

func getDefault(req *http.Request, transport *http.Transport) (*ResponseWithHeader, error) {
	return getUpstream(req, transport)
}

func getUpstream(req *http.Request, transport *http.Transport) (*ResponseWithHeader, error) {
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &ResponseWithHeader{
		InProgress: false,
		Expires:    time.Now().Add(time.Duration(randInt(1000, 1500)) * time.Second),
		Body:       body,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}, nil
}
