package logging

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/chubin/wttr.in/internal/util"
)

// Logging request.
//

// RequestLogger logs all incoming HTTP requests.
type RequestLogger struct {
	buf      map[logEntry]int
	filename string
	m        sync.Mutex

	period    time.Duration
	lastFlush time.Time
}

type logEntry struct {
	Proto     string
	IP        string
	URI       string
	UserAgent string
}

// NewRequestLogger returns a new RequestLogger for the specified log file.
// Flush logging entries after period of time.
//
// If filename is empty, no log will be written, and all logging entries
// will be silently dropped.
func NewRequestLogger(filename string, period time.Duration) *RequestLogger {
	return &RequestLogger{
		buf:      map[logEntry]int{},
		filename: filename,
		m:        sync.Mutex{},
		period:   period,
	}
}

// Log logs information about a HTTP request.
func (rl *RequestLogger) Log(r *http.Request) error {
	le := logEntry{
		Proto:     "http",
		IP:        util.ReadUserIP(r),
		URI:       r.RequestURI,
		UserAgent: r.Header.Get("User-Agent"),
	}
	if r.TLS != nil {
		le.Proto = "https"
	}

	// Do not log 127.0.0.1 connections
	if le.IP == "127.0.0.1" {
		return nil
	}

	rl.m.Lock()
	rl.buf[le]++
	rl.m.Unlock()

	if time.Since(rl.lastFlush) > rl.period {
		return rl.flush()
	}

	return nil
}

// flush stores log data to disk, and flushes the buffer.
func (rl *RequestLogger) flush() error {
	rl.m.Lock()
	defer rl.m.Unlock()

	// It is possible, that while waiting the mutex,
	// the buffer was already flushed.
	if time.Since(rl.lastFlush) <= rl.period {
		return nil
	}

	if rl.filename != "" {
		// Generate log output.
		output := ""
		for k, hitsNumber := range rl.buf {
			output += fmt.Sprintf("%s %3d %s\n", time.Now().Format(time.RFC3339), hitsNumber, k.String())
		}

		// Open log file.
		//nolint:nosnakecase
		f, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			return err
		}
		defer f.Close()

		// Save output to log file.
		_, err = f.Write([]byte(output))
		if err != nil {
			return err
		}
	}

	// Flush buffer.
	rl.buf = map[logEntry]int{}
	rl.lastFlush = time.Now()

	return nil
}

// String returns string representation of logEntry.
func (e *logEntry) String() string {
	return fmt.Sprintf(
		"%s %s %s %s",
		e.Proto,
		e.IP,
		e.URI,
		e.UserAgent,
	)
}
