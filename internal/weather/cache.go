// Package weather defines domain-level types and interfaces used across the application,
// including caching abstractions.
package weather

import (
	"net/http"
	"time"
)

// CacheEntry represents a cached HTTP response.
// It is immutable once stored in the cache.
type CacheEntry struct {
	// Body is the response body bytes
	Body []byte

	// Header contains the HTTP response headers to return to the client
	Header http.Header

	// StatusCode is the HTTP status code (200, 404, etc.)
	StatusCode int

	// Expires is the absolute time after which this entry should be considered stale
	Expires time.Time
}

// Cacher defines the contract for all cache implementations used by the request processor.
// Implementations may use LRU, Redis, in-memory maps, etc.
type Cacher interface {
	// Get returns a valid (non-expired) cache entry for the given key.
	// Returns nil if the key does not exist or the entry has expired.
	Get(key string) *CacheEntry

	// Set stores a completed response in the cache under the given key.
	// This should clear any in-progress state for the same key.
	Set(key string, entry CacheEntry)

	// SetInProgress marks that a request for this key is currently being processed.
	// Used to prevent duplicate upstream requests (coalescing).
	SetInProgress(key string)

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
	WaitForCompletion(key string, maxWait time.Duration) (*CacheEntry, error)

	// Remove deletes the entry (and any in-progress marker) for the given key.
	// Typically called when the upstream response should not be cached
	// (e.g. status 500, 429, etc.).
	Remove(key string)

	// Close releases any resources held by the cache (connections, goroutines, etc.).
	// Optional for in-memory implementations, required for networked caches.
	Close() error
}

// buildCacheKey generates a cache signature very similar to original wttr.in logic
func buildCacheKey(r *http.Request) string {
	ua := r.Header.Get("User-Agent")
	host := r.Host
	uri := r.RequestURI
	ip := getClientIP(r)
	lang := r.Header.Get("Accept-Language")

	// You can keep it exactly like original:
	return ua + ":" + host + uri + ":" + ip + ":" + lang
}
