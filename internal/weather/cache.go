// Package weather defines domain-level types and interfaces used across the application,
// including caching abstractions.
package weather

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chubin/wttr.in/internal/domain"
)

// Cacher defines the contract for all cache implementations used by the request processor.
// Implementations may use LRU, Redis, in-memory maps, etc.
type Cacher interface {
	// Get returns a valid (non-expired) cache entry for the given key.
	// Returns nil if the key does not exist or the entry has expired.
	Get(key string) *domain.CacheEntry

	// Set stores a completed response in the cache under the given key.
	// This should clear any in-progress state for the same key.
	Set(key string, entry domain.CacheEntry)

	// SetInProgressIfNotExists marks that a request for this key is currently being processed.
	// Used to prevent duplicate upstream requests (coalescing).
	SetInProgressIfNotExists(key string) bool

	// IsInProgress returns true if the key is currently being processed by another goroutine.
	IsInProgress(key string) bool

	// WaitForCompletion blocks until the in-progress state for this key is cleared
	// (i.e., the first request has finished and Set() has been called),
	// or until the timeout is reached.
	//
	// Returns:
	//   - the completed entry if available
	//   - nil, nil if the entry was removed (e.g. upstream error)
	//   - nil, error on timeout
	WaitForCompletion(key string, maxWait time.Duration) (*domain.CacheEntry, error)

	// Remove deletes the entry (and any in-progress marker) for the given key.
	// Typically called when the upstream response should not be cached
	// (e.g. status 500, 429, etc.).
	Remove(key string)

	// Close releases any resources held by the cache (connections, goroutines, etc.).
	// Optional for in-memory implementations, required for networked caches.
	Close() error
}

// DefaultCacheInterval defines how often the cache key "rotates" (i.e. effective TTL).
// All entries within the same interval share the same key.
const DefaultCacheInterval = 1 * time.Hour

// buildCacheKey generates a cache signature very similar to original wttr.in logic
func buildCacheKey(r *http.Request) string {
	ua := r.Header.Get("User-Agent")
	host := r.Host
	uri := r.RequestURI
	ip := getClientIP(r)
	lang := r.Header.Get("Accept-Language")

	// Bucket the current time by the interval
	// Example: for 1 hour interval, all requests in the same hour get the same bucket
	bucket := time.Now().Unix() / int64(DefaultCacheInterval.Seconds())

	return fmt.Sprintf("%s:%s%s:%s:%s:%d", ua, host, uri, ip, lang, bucket)
}
