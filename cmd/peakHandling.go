package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron"
)

func (rp *RequestProcessor) startPeakHandling() {
	c := cron.New()
	// cronTime := fmt.Sprintf("%d,%d * * * *", 30-prefetchInterval/60, 60-prefetchInterval/60)
	c.AddFunc(
		"24 * * * *",
		func() { rp.prefetchPeakRequests(&rp.peakRequest30) },
	)
	c.AddFunc(
		"54 * * * *",
		func() { rp.prefetchPeakRequests(&rp.peakRequest60) },
	)
	c.Start()
}

func (rp *RequestProcessor) savePeakRequest(cacheDigest string, r *http.Request) {
	_, min, _ := time.Now().Clock()
	if min == 30 {
		rp.peakRequest30.Store(cacheDigest, *r)
	} else if min == 0 {
		rp.peakRequest60.Store(cacheDigest, *r)
	}
}

func (rp *RequestProcessor) prefetchRequest(r *http.Request) {
	rp.ProcessRequest(r)
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

func (rp *RequestProcessor) prefetchPeakRequests(peakRequestMap *sync.Map) {
	peakRequestLen := syncMapLen(peakRequestMap)
	if peakRequestLen == 0 {
		return
	}
	log.Printf("PREFETCH: Prefetching %d requests\n", peakRequestLen)
	sleepBetweenRequests := time.Duration(rp.config.Uplink.PrefetchInterval*1000/peakRequestLen) * time.Millisecond
	peakRequestMap.Range(func(key interface{}, value interface{}) bool {
		go func(r http.Request) {
			rp.prefetchRequest(&r)
		}(value.(http.Request))
		peakRequestMap.Delete(key)
		time.Sleep(sleepBetweenRequests)
		return true
	})
}
