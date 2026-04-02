// Package cache provides a concrete implementation of the weather.Cacher
// interface using hashicorp/golang-lru/v2 as the underlying storage.
package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/weather"
)

const (
	inProgressMarker    = "__IN_PROGRESS__"
	defaultPollInterval = 25 * time.Millisecond
	defaultMaxWait      = 12 * time.Second
)

// LRUCacher is a Cacher implementation backed by an LRU cache with
// support for in-progress markers to prevent thundering herd problems.
type LRUCacher struct {
	cache        *lru.Cache[string, any]
	mu           sync.Mutex
	inProgress   map[string]struct{}
	pollInterval time.Duration
	maxWait      time.Duration
}

// NewLRU creates and returns a new LRU-based Cacher.
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

// Get returns a valid, non-expired cache entry or nil if not found/expired.
func (c *LRUCacher) Get(key string) *domain.CacheEntry {
	raw, ok := c.cache.Get(key)
	if !ok {
		return nil
	}

	// Check for in-progress marker
	if marker, ok := raw.(string); ok && marker == inProgressMarker {
		return nil
	}

	entry, ok := raw.(domain.CacheEntry)
	if !ok {
		// Invalid entry type → treat as miss
		c.cache.Remove(key)
		return nil
	}

	if time.Now().After(entry.Expires) {
		c.cache.Remove(key)
		return nil
	}

	return &entry
}

// Set stores a completed cache entry and clears the in-progress state.
func (c *LRUCacher) Set(key string, entry domain.CacheEntry) {
	c.mu.Lock()
	delete(c.inProgress, key)
	c.mu.Unlock()

	c.cache.Add(key, entry)
}

// SetInProgress marks a key as currently being processed to prevent duplicate work.
func (c *LRUCacher) SetInProgress(key string) {
	c.mu.Lock()
	c.inProgress[key] = struct{}{}
	c.mu.Unlock()

	c.cache.Add(key, inProgressMarker)
}

// IsInProgress checks whether the key is currently being processed.
func (c *LRUCacher) IsInProgress(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.inProgress[key]
	return exists
}

// WaitForCompletion blocks until the in-progress flag is cleared or timeout occurs.
// Returns the entry if available, nil if missing or upstream failed, or error on timeout.
func (c *LRUCacher) WaitForCompletion(key string, maxWait time.Duration) (*domain.CacheEntry, error) {
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
			// Entry was removed (likely upstream failure)
			return nil, nil
		}

		select {
		case <-ticker.C:
			// continue polling
		case <-time.After(time.Until(deadline)): // safety for edge case
			return nil, errors.New("timeout waiting for cache entry to complete")
		}
	}

	return nil, errors.New("timeout waiting for cache entry to complete")
}

// Remove deletes the entry for the given key (including any in-progress marker).
func (c *LRUCacher) Remove(key string) {
	c.mu.Lock()
	delete(c.inProgress, key)
	c.mu.Unlock()

	c.cache.Remove(key)
}

// Close performs cleanup.
func (c *LRUCacher) Close() error {
	c.mu.Lock()
	c.inProgress = nil
	c.mu.Unlock()

	// Optional: c.cache.Purge()
	return nil
}
