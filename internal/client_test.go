package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestConnectRequiresSettingsAndDefersHTTP(t *testing.T) {
	settings := connector.ConnectionConfig{Settings: map[string]string{baseURLKey: "https://dns.example.test"}}
	if _, err := (&Provider{}).Connect(settings); !errors.Is(err, connector.ErrUnavailable) {
		t.Fatalf("Connect() without HTTP = %v, want unavailable", err)
	}

	provider := Provider{HTTP: func(context.Context, abi.HTTPRequest) (abi.HTTPResponse, error) {
		t.Fatal("Connect must not make an HTTP request")
		return abi.HTTPResponse{}, nil
	}}
	if _, err := provider.Connect(connector.ConnectionConfig{}); !errors.Is(err, connector.ErrInvalidServerURL) {
		t.Fatalf("Connect() without endpoint = %v, want invalid URL", err)
	}
	if client, err := provider.Connect(settings); err != nil || client == nil {
		t.Fatalf("Connect() = %v, %v", client, err)
	}
}
