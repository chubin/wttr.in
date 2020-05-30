package main

import (
	"log"
	"sync"
	"time"
)

type safeCounter struct {
	v   map[int]int
	mux sync.Mutex
}

func (c *safeCounter) inc(key int) {
	c.mux.Lock()
	c.v[key]++
	c.mux.Unlock()
}

// func (c *safeCounter) val(key int) int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
// 	return c.v[key]
// }
//
// func (c *safeCounter) reset(key int) int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
// 	result := c.v[key]
// 	c.v[key] = 0
// 	return result
// }

var queriesPerMinute safeCounter

func printStat() {
	_, min, _ := time.Now().Clock()
	queriesPerMinute.inc(min)
	log.Printf("Processed %d requests\n", min)
}
