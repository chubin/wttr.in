package logging

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/chubin/wttr.in/internal/util"
)

// RequestLogger logs incoming HTTP requests with aggregation and periodic flushing.
type RequestLogger struct {
	buf       map[logEntry]int
	filename  string
	period    time.Duration
	lastFlush time.Time
	m         sync.Mutex
}

// logEntry represents a unique request signature for aggregation.
type logEntry struct {
	Proto     string
	IP        string
	URI       string
	UserAgent string
}

func (e logEntry) String() string {
	return fmt.Sprintf("%s %s %s %s", e.Proto, e.IP, e.URI, e.UserAgent)
}

// NewRequestLogger creates a new request logger.
// If filename is empty, logs are dropped silently.
// period must be > 0, otherwise it defaults to 1 minute.
func NewRequestLogger(filename string, period time.Duration) *RequestLogger {
	if period <= 0 {
		period = time.Minute
	}

	return &RequestLogger{
		buf:       make(map[logEntry]int),
		filename:  filename,
		period:    period,
		lastFlush: time.Now(),
	}
}

// Log records a request. Requests from 127.0.0.1 are ignored.
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

	// Skip localhost requests (common for health checks, metrics, etc.)
	if le.IP == "127.0.0.1" || le.IP == "::1" {
		return nil
	}

	rl.m.Lock()
	rl.buf[le]++
	shouldFlush := time.Since(rl.lastFlush) > rl.period
	rl.m.Unlock()

	if shouldFlush {
		return rl.flush()
	}
	return nil
}

// flush writes the current buffer to disk and resets it.
// It is safe to call concurrently.
func (rl *RequestLogger) flush() error {
	rl.m.Lock()
	defer rl.m.Unlock()

	// Double-check after acquiring lock
	if time.Since(rl.lastFlush) <= rl.period {
		return nil
	}

	if len(rl.buf) == 0 {
		rl.lastFlush = time.Now()
		return nil
	}

	if rl.filename != "" {
		output := rl.buildLogOutput()

		f, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", rl.filename, err)
		}
		defer f.Close()

		if _, err := f.WriteString(output); err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
	}

	// Only clear buffer and update timestamp if write succeeded (or no file)
	rl.buf = make(map[logEntry]int)
	rl.lastFlush = time.Now()

	return nil
}

// buildLogOutput generates the log lines from current buffer.
func (rl *RequestLogger) buildLogOutput() string {
	var sb []byte
	now := time.Now().Format(time.RFC3339)

	for k, hits := range rl.buf {
		line := fmt.Sprintf("%s %3d %s\n", now, hits, k.String())
		sb = append(sb, line...)
	}
	return string(sb)
}

// Flush forces a flush of the current buffer. Useful for graceful shutdown.
func (rl *RequestLogger) Flush() error {
	return rl.flush()
}