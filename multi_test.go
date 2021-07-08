package apibuilder

import (
	"encoding/json"
	"testing"
)

func TestMulti_AddRoute(t *testing.T) {
	rr, err := testResponse(createRoute().handlerFunc)
	if err != nil {
		t.Fatal(rr)
	}
	var calls []string
	if err := json.NewDecoder(rr.Body).Decode(&calls); err != nil {
		t.Fatalf("could not decode response")
	}
	if len(calls) != 2 { //nolint:gocritic
		t.Fatalf("expected stack call of len 2: got %d", len(calls))
	}
	if calls[0] != "1" {
		t.Fatalf("expected first call of val 1, got %s", calls[0])
	}
}

func createRoute() API {
	return NewMulti().API("GET;POST;PUT;PATCH", "/api/v1/test", simpleHandler, stackWithValue("1"), stackWithValue("2")).Done()[0]
}
