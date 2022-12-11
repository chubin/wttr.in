package types

import "errors"

type CacheType string

const (
	CacheTypeDB    = "db"
	CacheTypeFiles = "files"
)

var (
	ErrNotFound          = errors.New("cache entry not found")
	ErrInvalidCacheEntry = errors.New("invalid cache entry format")
)
