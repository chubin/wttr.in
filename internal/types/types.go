package types

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrUnknownLocationService = errors.New("unknown location service")
	ErrUpstream               = errors.New("upstream error")
	ErrInvalidCacheEntry      = errors.New("invalid cache entry")
)

// NotFoundError is a user-facing "unknown location" error carrying a
// localized message. It is rendered as an HTTP 404 response instead of
// the generic 500 error page.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }

type Cadre struct {
	Body []byte
}
