//go:build integration

package internal

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestProbeAgainstExampleUpstream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/health" || r.Header.Get("Accept") != "application/json" {
			http.NotFound(w, r)
			return
		}
		_, _ = io.WriteString(w, `{"version":"1.2.3"}`)
	}))
	t.Cleanup(upstream.Close)

	// Stand in for Izena's HTTP host while exercising a real fixture endpoint.
	host := func(ctx context.Context, call abi.HTTPRequest) (abi.HTTPResponse, error) {
		request, err := http.NewRequestWithContext(ctx, call.Method, upstream.URL+"/"+strings.Join(call.Parts, "/"), nil)
		if err != nil {
			return abi.HTTPResponse{}, err
		}
		request.Header = call.Headers.Clone()
		response, err := upstream.Client().Do(request)
		if err != nil {
			return abi.HTTPResponse{}, err
		}
		defer func() { _ = response.Body.Close() }()
		body, err := io.ReadAll(io.LimitReader(response.Body, int64(call.Limit)))
		return abi.HTTPResponse{Status: response.StatusCode, Body: body}, err
	}
	result, err := Dispatch(t.Context(), Provider{HTTP: host}, abi.Request{
		Operation: abi.Probe, Settings: map[string]string{baseURLKey: upstream.URL},
	})
	want := connector.Profile{DaemonType: providerName, Version: "1.2.3"}
	if err != nil || result != want {
		t.Fatalf("probe result = %+v, %v, want %+v", result, err, want)
	}
}
