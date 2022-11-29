package routing

import (
	"net/http"
	"time"
)

// CadreFormat specifies how the shot data is formatted.
type CadreFormat int

const (
	// CadreFormatANSI represents Terminal ANSI format.
	CadreFormatANSI = iota

	// CadreFormatHTML represents HTML.
	CadreFormatHTML

	// CadreFormatPNG represents PNG.
	CadreFormatPNG
)

// Cadre contains result of a query execution.
type Cadre struct {
	// Body contains the data of Cadre, formatted as Format.
	Body []byte

	// Format of the shot.
	Format CadreFormat

	// Expires contains the time of the Cadre expiration,
	// or 0 if it does not expire.
	Expires time.Time
}

// Handler can handle queries and return views.
type Handler interface {
	Response(*http.Request) *Cadre
}

type routeFunc func(*http.Request) bool

type route struct {
	routeFunc
	Handler
}

// Router keeps a routing table, and finds queries handlers, based on its rules.
type Router struct {
	rt []route
}

// Route returns a query handler based on its content.
func (r *Router) Route(req *http.Request) Handler {
	for _, re := range r.rt {
		if re.routeFunc(req) {
			return re.Handler
		}
	}
	return nil
}

// AddPath adds route for a static path.
func (r *Router) AddPath(path string, handler Handler) {
	r.rt = append(r.rt, route{routePath(path), handler})
}

func routePath(path string) routeFunc {
	return routeFunc(func(req *http.Request) bool {
		return req.URL.Path == path
	})
}
