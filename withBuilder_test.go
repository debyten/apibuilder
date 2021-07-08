package apibuilder

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type stackKey int

const stackCall stackKey = iota

func Test_ApiBuilding(t *testing.T) {
	NewMulti().With(preflightRequest, basicAuth).
		API("GET", "/api/test", simpleHandler)
}

func preflightRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r)
	}
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if u != "admin" || p != "admin" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func TestMultiBuilder_With(t *testing.T) {

	h := createHandler().handlerFunc
	rr, err := testResponse(h)
	if err != nil {
		t.Fatal(err)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
	var calls []string
	if err := json.NewDecoder(rr.Body).Decode(&calls); err != nil {
		t.Fatalf("could not decode response")
	}
	if len(calls) != 2 { //nolint:gocritic
		t.Fatalf("expected stack call of len 2: got %d", len(calls))
	}
	for i, call := range calls {
		if call != strconv.Itoa(i) {
			t.Fatalf("expected stack call %d value '%d', got %s", i, i, call)
		}
	}
}

func TestMultiBuilderMiddleware_With(t *testing.T) {
	h := createHandlerWithMiddleware().handlerFunc
	rr, err := testResponse(h)
	if err != nil {
		t.Fatal(err)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
	var calls []string
	if err := json.NewDecoder(rr.Body).Decode(&calls); err != nil {
		t.Fatalf("could not decode response")
	}
	if len(calls) != 4 { //nolint:gocritic
		t.Fatalf("expected stack call of len 4: got %d", len(calls))
	}
	for i, call := range calls {
		if call != strconv.Itoa(i) {
			t.Fatalf("expected stack call %d value '%d', got %s", i, i, call)
		}
	}
}

func testResponse(fn func(w http.ResponseWriter, r *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "/testMiddleware", nil)
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	fn(rr, req)
	return rr, nil
}

func createHandler() API {
	return NewMulti().With(stackWithValue("0"), stackWithValue("1")).
		API("GET", "/testMiddleware", simpleHandler).End().Done()[0]
}

func createHandlerWithMiddleware() API {
	return NewMulti().With(stackWithValue("0"), stackWithValue("1")).
		API("GET", "/testMiddleware", simpleHandler, stackWithValue("2"), stackWithValue("3")).End().Done()[0]
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	stack, ok := r.Context().Value(stackCall).([]string)
	if !ok {
		w.WriteHeader(500)
		return
	}
	if err := json.NewEncoder(w).Encode(&stack); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func stackWithValue(val string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			stack, ok := r.Context().Value(stackCall).([]string)
			if !ok {
				ctx := context.WithValue(r.Context(), stackCall, []string{val})
				h(w, r.WithContext(ctx))
				return
			}
			stack = append(stack, val)
			ctx := context.WithValue(r.Context(), stackCall, stack)
			h(w, r.WithContext(ctx))
		}
	}
}
