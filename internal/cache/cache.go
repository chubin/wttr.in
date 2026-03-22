// Package cache provides a concrete implementation of the weather.Cacher
// interface using hashicorp/golang-lru as the underlying storage.
package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/chubin/wttr.in/internal/weather"
)

const (
	inProgressMarker    = "__IN_PROGRESS__"
	defaultPollInterval = 25 * time.Millisecond
	defaultMaxWait      = 12 * time.Second
)

// LRUCacher is a Cacher implementation backed by an LRU cache.
type LRUCacher struct {
	cache        *lru.Cache[string, any]
	mu           sync.Mutex
	inProgress   map[string]struct{}
	pollInterval time.Duration
	maxWait      time.Duration
}

// NewLRU creates and returns a new Cacher using LRU eviction policy.
func NewLRU(cfg Config) (weather.Cacher, error) {
	size := cfg.Size

	if size <= 0 {
		size = 1024 // sensible default
	}

	lruCache, err := lru.New[string, any](size)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	return &LRUCacher{
		cache:        lruCache,
		inProgress:   make(map[string]struct{}, 512),
		pollInterval: defaultPollInterval,
		maxWait:      defaultMaxWait,
	}, nil
}

// Get returns a valid, non-expired cache entry or nil.
func (c *LRUCacher) Get(key string) *weather.CacheEntry {
	raw, ok := c.cache.Get(key)
	if !ok {
		return nil
	}

	// Check for in-progress marker
	if marker, isString := raw.(string); isString && marker == inProgressMarker {
		return nil
	}

	entry, ok := raw.(weather.CacheEntry)
	if !ok {
		return nil
	}

	if time.Now().After(entry.Expires) {
		c.cache.Remove(key)
		return nil
	}

	return &entry
}

// Set stores a completed cache entry and clears the in-progress state.
func (c *LRUCacher) Set(key string, entry weather.CacheEntry) {
	c.mu.Lock()
	delete(c.inProgress, key)
	c.mu.Unlock()

	c.cache.Add(key, entry)
}

// SetInProgress marks a key as currently being processed.
func (c *LRUCacher) SetInProgress(key string) {
	c.mu.Lock()
	c.inProgress[key] = struct{}{}
	c.mu.Unlock()

	c.cache.Add(key, inProgressMarker)
}

// IsInProgress checks whether the key is currently being processed.
func (c *LRUCacher) IsInProgress(key string) bool {
	c.mu.Lock()
	_, exists := c.inProgress[key]
	c.mu.Unlock()
	return exists
}

// WaitForCompletion blocks until the in-progress flag is cleared or timeout occurs.
func (c *LRUCacher) WaitForCompletion(key string, maxWait time.Duration) (*weather.CacheEntry, error) {
	if maxWait <= 0 {
		maxWait = c.maxWait
	}

	deadline := time.Now().Add(maxWait)
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		if !c.IsInProgress(key) {
			if entry := c.Get(key); entry != nil {
				return entry, nil
			}
			// Entry was removed → upstream likely failed
			return nil, nil
		}

		<-ticker.C
	}

	return nil, errors.New("timeout waiting for cache entry to complete")
}

// Remove deletes the entry for the given key (including in-progress marker).
func (c *LRUCacher) Remove(key string) {
	c.mu.Lock()
	delete(c.inProgress, key)
	c.mu.Unlock()
	c.cache.Remove(key)
}

// Close performs any necessary cleanup (optional for LRU).
func (c *LRUCacher) Close() error {
	c.mu.Lock()
	c.inProgress = nil
	c.mu.Unlock()
	// LRU doesn't require explicit close, but we could purge if desired
	// c.cache.Purge()
	return nil
}
