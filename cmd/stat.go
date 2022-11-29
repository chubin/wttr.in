package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/chubin/wttr.in/internal/routing"
)

// Stats holds processed requests statistics.
type Stats struct {
	m         sync.Mutex
	v         map[string]int
	startTime time.Time
}

// NewStats returns new Stats.
func NewStats() *Stats {
	return &Stats{
		v:         map[string]int{},
		startTime: time.Now(),
	}
}

// Inc key by one.
func (c *Stats) Inc(key string) {
	c.m.Lock()
	c.v[key]++
	c.m.Unlock()
}

// Get current key counter value.
func (c *Stats) Get(key string) int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.v[key]
}

// Reset key counter.
func (c *Stats) Reset(key string) int {
	c.m.Lock()
	defer c.m.Unlock()
	result := c.v[key]
	c.v[key] = 0
	return result
}

// Show returns current statistics formatted as []byte.
func (c *Stats) Show() []byte {
	var (
		b bytes.Buffer
	)

	c.m.Lock()
	defer c.m.Unlock()

	uptime := time.Since(c.startTime) / time.Second

	fmt.Fprintf(&b, "%-20s: %v\n", "Running since", c.startTime.Format(time.RFC3339))
	fmt.Fprintf(&b, "%-20s: %d\n", "Uptime (min)", uptime/60)

	fmt.Fprintf(&b, "%-20s: %d\n", "Total queries", c.v["total"])
	if uptime != 0 {
		fmt.Fprintf(&b, "%-20s: %d\n", "Throughput (QpM)", c.v["total"]*60/int(uptime))
	}

	fmt.Fprintf(&b, "%-20s: %d\n", "Cache L1 queries", c.v["cache1"])
	if c.v["total"] != 0 {
		fmt.Fprintf(&b, "%-20s: %d\n", "Cache L1 queries (%)", (100*c.v["cache1"])/c.v["total"])
	}

	fmt.Fprintf(&b, "%-20s: %d\n", "Upstream queries", c.v["total"]-c.v["cache1"])
	fmt.Fprintf(&b, "%-20s: %d\n", "Queries with format", c.v["format"])
	fmt.Fprintf(&b, "%-20s: %d\n", "Queries with format=j1", c.v["format=j1"])

	return b.Bytes()
}

func (c *Stats) Response(*http.Request) *routing.Cadre {
	return &routing.Cadre{
		Body: c.Show(),
	}
}
