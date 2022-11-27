package main

import "net/http"

type Handler interface {
	Response(*http.Request) *ResponseWithHeader
}

type routeFunc func(*http.Request) bool

type route struct {
	routeFunc
	Handler
}

type Router struct {
	rt []route
}

func (r *Router) Route(req *http.Request) Handler {
	for _, re := range r.rt {
		if re.routeFunc(req) {
			return re.Handler
		}
	}
	return nil
}

func (r *Router) AddPath(path string, handler Handler) {
	r.rt = append(r.rt, route{routePath(path), handler})
}

func routePath(path string) routeFunc {
	return routeFunc(func(req *http.Request) bool {
		return req.URL.Path == path
	})
}
