package apibuilder

import (
	"net/http"
	"strings"
)

// With inject parent [] Middleware into *Multi instance
type With struct {
	original   []Middleware
	middleware []Middleware
	m          *Multi
}

// With starts building subsequent routes with the specified middleware functions.
// The middleware functions will be executed in order.
//
// Example:
//
//	m.With(firstMiddleware, secondMiddleware,...).
//	  API("GET", "/path", handler1).
//	  API("POST", "/path2", handler2, myMiddlewareFunc).End()
//
// Result stack:
//
//	(GET /path): firstMiddleware => secondMiddleware => handler1
//	(GET /path2): firstMiddleware => secondMiddleware => myMiddlewareFunc => handler2
func (m *Multi) With(middleware ...Middleware) *With {
	return &With{
		original:   middleware,
		middleware: middleware,
		m:          m,
	}
}

func (w *With) NewGroup(middleware ...Middleware) *With {
	return &With{
		original:   w.middleware,
		middleware: append(w.original, middleware...),
		m:          w.m,
	}
}

// API calls Multi.API injecting previously specified With middlewares
func (w *With) API(method string, path string, h http.HandlerFunc, mid ...Middleware) *With {
	if strings.Contains(method, ";") {
		panic("please use the With.APIs() method")
	}
	h = concat(concat(h, mid), w.middleware)
	w.m.API(method, path, h)
	return w
}

// APIs calls Multi.API injecting previously specified With middlewares
func (w *With) APIs(methods []string, path string, h http.HandlerFunc, mid ...Middleware) *With {
	h = concat(concat(h, mid), w.middleware)
	w.m.API(strings.Join(methods, ";"), path, h)
	return w
}

// End returns the *Multi instance
func (w *With) End() *Multi {
	return w.m
}

// reverse is needed to preserve the order of execution of []Middleware
func reverse(s []Middleware) []Middleware {
	if s == nil {
		return nil
	}
	a := make([]Middleware, len(s))
	copy(a, s)
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}
