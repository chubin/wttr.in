package util

import (
	"log"
	"net"
	"net/http"
)

// ReadUserIP returns IP address of the client from http.Request,
// taking into account the HTTP headers.
func ReadUserIP(r *http.Request) string {
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
