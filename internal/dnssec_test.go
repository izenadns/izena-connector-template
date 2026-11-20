package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestReadOnlyConnectorRejectsDNSSECOperations(t *testing.T) {
	client := Client{}
	if _, err := client.GetDNSSEC(t.Context(), exampleZoneID); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("GetDNSSEC() = %v, want unsupported", err)
	}
	if err := client.SetDNSSEC(t.Context(), exampleZoneID, true); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("SetDNSSEC() = %v, want unsupported", err)
	}
}
