package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestReadOnlyConnectorRejectsZoneWrites(t *testing.T) {
	client := Client{}
	if _, err := client.CreateZone(t.Context(), connector.ZoneCreateInput{Name: "new.example.test."}); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("CreateZone() = %v, want unsupported", err)
	}
	got, err := client.UpdateZone(t.Context(), exampleZoneID, connector.ZoneUpdateInput{Kind: "example"})
	if !errors.Is(err, connector.ErrUnsupported) || got.Attempted {
		t.Fatalf("UpdateZone() = %+v, %v, want unsupported without an attempt", got, err)
	}
	if err := client.DeleteZone(t.Context(), exampleZoneID); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("DeleteZone() = %v, want unsupported", err)
	}
}
