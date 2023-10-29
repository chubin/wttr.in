package types

import "errors"

var (
	ErrNotFound          = errors.New("cache entry not found")
	ErrInvalidCacheEntry = errors.New("invalid cache entry format")
	ErrUpstream          = errors.New("upstream error")

	// ErrNoServersConfigured means that there are no servers to run.
	ErrNoServersConfigured = errors.New("no servers configured")

	ErrUnknownLocationService = errors.New("unknown location service")
)
