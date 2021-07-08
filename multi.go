package apibuilder

import (
	"net/http"
)

// Multi represents a collection of built APIs
type Multi struct {
	routes []API
}

// NewMulti returns a new *Multi instance
func NewMulti() *Multi {
	return &Multi{routes: make([]API, 0)}
}

// API saves an API into *Multi instance
//
// Example:
//  multi := NewMulti()
//  multi.API("GET", "/api/v1/users", func..., m1, m2, m3)
func (m *Multi) API(method string, path string, h http.HandlerFunc, mid ...Middleware) *Multi {
	cb := New(method).Path(path).Handler(h)
	if len(mid) > 0 {
		cb = cb.Middleware(mid...)
	}
	m.routes = append(m.routes, cb.Build())
	return m
}

// Done returns the slice of API built
func (m *Multi) Done() []API {
	return m.routes
}
