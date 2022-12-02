package config

// Config of the program.
type Config struct {
	Cache
	Logging
	Server
	Uplink
}

// Logging configuration.
type Logging struct {

	// AccessLog path.
	AccessLog string

	// ErrorsLog path.
	ErrorsLog string

	// Interval between access log flushes, in seconds.
	Interval int
}

// Server configuration.
type Server struct {

	// PortHTTP is port where HTTP server must listen.
	// If 0, HTTP is disabled.
	PortHTTP int

	// PortHTTP is port where the HTTPS server must listen.
	// If 0, HTTPS is disabled.
	PortHTTPS int

	// TLSCertFile contains path to cert file for TLS Server.
	TLSCertFile string

	// TLSCertFile contains path to key file for TLS Server.
	TLSKeyFile string
}

// Uplink configuration.
type Uplink struct {
	// Address contains address of the uplink server in form IP:PORT.
	Address string

	// Timeout for upstream queries.
	Timeout int

	// PrefetchInterval contains time (in milliseconds) indicating,
	// how long the prefetch procedure should take.
	PrefetchInterval int
}

// Cache configuration.
type Cache struct {
	// Size of the main cache.
	Size int
}

// Default contains the default configuration.
func Default() *Config {
	return &Config{
		Cache{
			Size: 12800,
		},
		Logging{
			AccessLog: "/wttr.in/log/access.log",
			ErrorsLog: "/wttr.in/log/errors.log",
			Interval:  300,
		},
		Server{
			PortHTTP:    8083,
			PortHTTPS:   8084,
			TLSCertFile: "/wttr.in/etc/fullchain.pem",
			TLSKeyFile:  "/wttr.in/etc/privkey.pem",
		},
		Uplink{
			Address:          "127.0.0.1:9002",
			Timeout:          30,
			PrefetchInterval: 300,
		},
	}
}

// Conf contains the current configuration
var Conf = Default()
