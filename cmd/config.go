package main

// Config of the program.
type Config struct {
	Logging
	Server
}

// Logging configuration.
type Logging struct {

	// AccessLog path.
	AccessLog string

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

// Conf contains the current configuration.
var Conf = Config{
	Logging{
		AccessLog: "/wttr.in/log/access.log",
		Interval:  300,
	},
	Server{
		PortHTTP:    8083,
		PortHTTPS:   8084,
		TLSCertFile: "/wttr.in/etc/fullchain.pem",
		TLSKeyFile:  "/wttr.in/etc/privkey.pem",
	},
}
