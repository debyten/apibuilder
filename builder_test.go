package apibuilder

import (
	"net/http"
	"testing"
)

func TestRoute_GetMethod(t *testing.T) {
	type fields struct {
		method      string
		path        string
		handlerFunc http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "t1", fields: fields{method: "GET", path: "/api/test", handlerFunc: nil}, want: "GET"},
	}
	for _, tt := range tests {
		ts := tt
		t.Run(tt.name, func(t *testing.T) {
			r := &API{
				method:      ts.fields.method,
				path:        ts.fields.path,
				handlerFunc: ts.fields.handlerFunc,
			}
			if got := r.Method(); got != ts.want {
				t.Errorf("API.Method() = %v, want %v", got, ts.want)
			}
		})
	}
}

func TestRoute_GetPath(t *testing.T) {
	type fields struct {
		method      string
		path        string
		handlerFunc http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "t1", fields: fields{method: "GET", path: "/api/test", handlerFunc: nil}, want: "/api/test"},
	}
	for _, tt := range tests {
		ts := tt
		t.Run(tt.name, func(t *testing.T) {
			r := &API{
				method:      ts.fields.method,
				path:        ts.fields.path,
				handlerFunc: ts.fields.handlerFunc,
			}
			if got := r.Path(); got != ts.want {
				t.Errorf("API.Path() = %v, want %v", got, ts.want)
			}
		})
	}
}

func TestBuilder_Build(t *testing.T) {
	r := New("GET").Handler(simpleHandler).Build()
	rr, err := testResponse(r.handlerFunc)
	if err != nil {
		t.Fatalf("got error on testresponse: %v", err)
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
