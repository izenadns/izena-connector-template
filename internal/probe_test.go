package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestProbeClassifiesUpstreamResponses(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{"authentication", http.StatusUnauthorized, "", connector.ErrAuthentication},
		{"upstream unavailable", http.StatusServiceUnavailable, "", connector.ErrUnavailable},
		{"malformed response", http.StatusOK, "{", connector.ErrMalformed},
		{"missing version", http.StatusOK, "{}", connector.ErrMalformed},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := Client{exchange: func(context.Context, abi.HTTPRequest) (abi.HTTPResponse, error) {
				return abi.HTTPResponse{Status: test.status, Body: []byte(test.body)}, nil
			}}
			if _, err := client.Probe(t.Context()); !errors.Is(err, test.want) {
				t.Fatalf("Probe() = %v, want %v", err, test.want)
			}
		})
	}
}
