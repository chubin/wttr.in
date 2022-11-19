package main

// Config of the program.
type Config struct {
	Logging
}

// Logging configuration.
type Logging struct {

	// AccessLog path.
	AccessLog string

	// Interval between access log flushes, in seconds.
	Interval int
}

// Conf contains the current configuration.
var Conf = Config{
	Logging{
		AccessLog: "/wttr.in/log/access.log",
		Interval:  300,
	},
}
