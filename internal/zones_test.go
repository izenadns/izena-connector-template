package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestZoneReadsUseOpaqueIDs(t *testing.T) {
	client := Client{}
	zones, err := client.ListZones(t.Context())
	if err != nil || len(zones) != 1 || zones[0].ID != exampleZoneID {
		t.Fatalf("ListZones() = %+v, %v", zones, err)
	}
	if _, err := client.GetZone(t.Context(), exampleZone); !errors.Is(err, connector.ErrNotFound) {
		t.Fatalf("GetZone() by name = %v, want not found", err)
	}
	if _, err := client.GetZoneMetadata(t.Context(), "missing"); !errors.Is(err, connector.ErrNotFound) {
		t.Fatalf("GetZoneMetadata() = %v, want not found", err)
	}
}
