package types

import "errors"

var (
	ErrNotFound          = errors.New("cache entry not found")
	ErrInvalidCacheEntry = errors.New("invalid cache entry format")
	ErrUpstream          = errors.New("upstream error")
)
