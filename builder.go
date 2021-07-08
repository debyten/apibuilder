package apibuilder

import (
	"net/http"
	"strings"
)

// Middleware describe a middleware for http.HandlerFunc
//
// Example:
//  func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//       // basic authentication check
//       h(w, r)
//     }
//  }
type Middleware func(h http.HandlerFunc) http.HandlerFunc

// API main struct that describes a built route.
//
//Example:
//  GET /api/v1/users func(http.ResponseWriter, *http.Request)
type API struct {
	method      string
	path        string
	handlerFunc http.HandlerFunc
}

// Handler describes a `func(http.ResponseWriter, *http.Request)`.
//
// The final Handler become the result of merged middlewares plus the handler itself.
//
// Example:
//  New("GET").Path("/").Handler(myHandler).Middleware(mid1, mid2, midN...)
// Result execution stack:
//  mid1 -> mid2 -> midN -> myHandler
func (a *API) Handler() http.HandlerFunc {
	return a.handlerFunc
}

// Path describe the api path
func (a *API) Path() string {
	return a.path
}

// Method describes the api method (GET, POST, PUT...)
func (a *API) Method() string {
	return a.method
}

// Methods describe the api methods when Method() is specified with the separator ";"
//
//Example:
//  New("GET;OPTIONS")
//Result:
//	[]string{"GET", "OPTIONS"}
func (a *API) Methods() []string {
	return strings.Split(a.method, ";")
}

// Builder represents a helper struct to build a single API
type Builder struct {
	method      string
	path        string
	handlerFunc http.HandlerFunc
	middleware  []Middleware
}

// New start building route specifying the method parameter (POST, PUT, PATCH...)
// **Note** you can specify multiple methods for a single route using the ';' separator (e.g. "GET;POST;PUT")
func New(method string) *Builder {
	return &Builder{method: method}
}

// Path describe the builder path, e.g. /api/v1/test
func (rb *Builder) Path(p string) *Builder {
	rb.path = p
	return rb
}

// Handler is a simple http.HandlerFunc
func (rb *Builder) Handler(h http.HandlerFunc) *Builder {
	rb.handlerFunc = h
	return rb
}

// Middleware sets the collection of middlewares in the Builder struct for a final concatenation with http.HandlerFunc in the Build method.
// the slice of middleware will be reversed to respect the sequentiality of the call stack
//
//Example:
//  Middleware(m1, m2, m3) => m1 -> m2 -> m3
//  Middleware(m7, m5, m9) => m7 -> m5 -> m9
func (rb *Builder) Middleware(middleware ...Middleware) *Builder {
	if len(rb.middleware) == 0 {
		rb.middleware = make([]Middleware, 0)
	}
	rb.middleware = append(rb.middleware, middleware...)
	return rb
}

// Build finalize the api build process applying a reverse function on middleware slice (to preserve the order) and
// returns the API with the handler func concatenated with the middlewares
func (rb *Builder) Build() API {
	rb.handlerFunc = concat(rb.handlerFunc, rb.middleware)
	return API{rb.method, rb.path, rb.handlerFunc}
}

func concat(h http.HandlerFunc, middleware []Middleware) http.HandlerFunc {
	for _, m := range reverse(middleware) {
		h = m(h)
	}
	return h
}
