package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestPrepareSettingsCanonicalizesEndpoint(t *testing.T) {
	p := Provider{}
	raw := map[string]string{baseURLKey: " https://dns.example.test/api/ "}
	got, err := p.PrepareSettings(t.Context(), raw)
	if err != nil || got.TargetKey != "https://dns.example.test/api" || raw[baseURLKey] != " https://dns.example.test/api/ " {
		t.Fatalf("PrepareSettings() = %+v, %v", got, err)
	}
	for _, endpoint := range []string{"", "/relative", "ftp://dns.example.test", "https://dns.example.test/?q=1"} {
		if _, err := p.PrepareSettings(t.Context(), map[string]string{baseURLKey: endpoint}); !errors.Is(err, connector.ErrInvalidServerURL) {
			t.Fatalf("PrepareSettings(%q) = %v, want invalid URL", endpoint, err)
		}
	}
}
