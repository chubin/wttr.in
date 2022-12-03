package logging

import (
	"os"
	"strings"
	"sync"
)

// LogSuppressor provides io.Writer interface for logging
// with lines suppression. For usage with log.Logger.
type LogSuppressor struct {
	filename   string
	suppress   []string
	linePrefix string

	logFile *os.File
	m       sync.Mutex
}

// NewLogSuppressor creates a new LogSuppressor for specified
// filename and lines to be suppressed.
func NewLogSuppressor(filename string, suppress []string, linePrefix string) *LogSuppressor {
	return &LogSuppressor{
		filename:   filename,
		suppress:   suppress,
		linePrefix: linePrefix,
	}
}

// Open opens log file.
func (ls *LogSuppressor) Open() error {
	var err error
	ls.logFile, err = os.OpenFile(ls.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	return err
}

// Close closes log file.
func (ls *LogSuppressor) Close() error {
	return ls.logFile.Close()
}

// Write writes p to log, and returns number f bytes written.
// Implements io.Writer interface.
func (ls *LogSuppressor) Write(p []byte) (n int, err error) {
	var (
		output string
	)

	ls.m.Lock()
	defer ls.m.Unlock()

	lines := strings.Split(string(p), ls.linePrefix)
	for _, line := range lines {
		if (func() bool {
			for _, suppress := range ls.suppress {
				if strings.Contains(line, suppress) {
					return true
				}
			}
			return false
		})() {
			continue
		}
		output += line
	}

	n, err = ls.logFile.Write([]byte(output))
	return n, err
}
