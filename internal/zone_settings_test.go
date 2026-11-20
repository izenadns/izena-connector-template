package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestZoneSettingsAreUnsupported(t *testing.T) {
	client := &Client{}
	if _, err := client.GetZoneSettings(t.Context(), exampleZoneID); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("GetZoneSettings() = %v, want unsupported", err)
	}
	if err := client.PatchZoneSettings(t.Context(), exampleZoneID, connector.ZoneSettingsPatch{Reset: []string{"serial_policy"}}); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("PatchZoneSettings() = %v, want unsupported", err)
	}
}
