package types

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrUnknownLocationService = errors.New("unknown location service")
	ErrUpstream               = errors.New("upstream error")
	ErrInvalidCacheEntry      = errors.New("invalid cache entry")
)

type Cadre struct {
	Body []byte
}
