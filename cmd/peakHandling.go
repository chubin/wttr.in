package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron"
)

var peakRequest30 sync.Map
var peakRequest60 sync.Map

func initPeakHandling() {
	c := cron.New()
	// cronTime := fmt.Sprintf("%d,%d * * * *", 30-prefetchInterval/60, 60-prefetchInterval/60)
	c.AddFunc("24 * * * *", prefetchPeakRequests30)
	c.AddFunc("54 * * * *", prefetchPeakRequests60)
	c.Start()
}

func savePeakRequest(cacheDigest string, r *http.Request) {
	_, min, _ := time.Now().Clock()
	if min == 30 {
		peakRequest30.Store(cacheDigest, *r)
	} else if min == 0 {
		peakRequest60.Store(cacheDigest, *r)
	}
}

func prefetchRequest(r *http.Request) {
	processRequest(r)
}

func syncMapLen(sm *sync.Map) int {
	count := 0

	f := func(key, value interface{}) bool {

		// Not really certain about this part, don't know for sure
		// if this is a good check for an entry's existence
		if key == "" {
			return false
		}
		count++

		return true
	}

	sm.Range(f)

	return count
}

func prefetchPeakRequests(peakRequestMap *sync.Map) {
	peakRequestLen := syncMapLen(peakRequestMap)
	log.Printf("PREFETCH: Prefetching %d requests\n", peakRequestLen)
	if peakRequestLen == 0 {
		return
	}
	sleepBetweenRequests := time.Duration(prefetchInterval*1000/peakRequestLen) * time.Millisecond
	peakRequestMap.Range(func(key interface{}, value interface{}) bool {
		r := value.(http.Request)
		log.Printf("Prefetching %s\n", key)
		prefetchRequest(&r)
		peakRequestMap.Delete(key)
		time.Sleep(sleepBetweenRequests)
		return true
	})
}

func prefetchPeakRequests30() {
	prefetchPeakRequests(&peakRequest30)
}

func prefetchPeakRequests60() {
	prefetchPeakRequests(&peakRequest60)
}
