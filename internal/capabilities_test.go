package internal

import (
	"slices"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestCapabilitiesAdvertiseOnlyExampleReadSupport(t *testing.T) {
	got, err := (&Client{}).Capabilities(t.Context())
	if err != nil || got.Provider != providerID || got.RecordModel != connector.RecordModelRRsets || !slices.Equal(got.RecordTypes, []string{"A"}) ||
		len(got.ZoneKinds) != 0 || len(got.Actions) != 0 || got.DNSSEC.Signing != connector.SigningUnsupported {
		t.Fatalf("Capabilities() = %+v, %v", got, err)
	}
}
