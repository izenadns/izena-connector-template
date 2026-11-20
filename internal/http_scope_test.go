package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestHTTPPolicyOnlyAllowsTheProbeRoute(t *testing.T) {
	calls := 0
	host := func(context.Context, abi.HTTPRequest) (abi.HTTPResponse, error) {
		calls++
		return abi.HTTPResponse{Status: http.StatusOK}, nil
	}
	probe := abi.HTTPRequest{Method: http.MethodGet, Parts: []string{"health"}}
	if _, err := scopeHTTP(abi.Probe, host)(t.Context(), probe); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		operation abi.Operation
		request   abi.HTTPRequest
	}{
		{abi.ListZones, probe},
		{abi.Probe, abi.HTTPRequest{Method: http.MethodPost, Parts: []string{"health"}}},
		{abi.Probe, abi.HTTPRequest{Method: http.MethodGet, Parts: []string{"health", "detail"}}},
	} {
		if _, err := scopeHTTP(test.operation, host)(t.Context(), test.request); !errors.Is(err, connector.ErrValidation) {
			t.Fatalf("unexpected scope result for %s: %v", test.operation, err)
		}
	}
	if calls != 1 {
		t.Fatalf("host calls = %d, want one permitted call", calls)
	}
}
