package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestReadOnlyConnectorRejectsZoneActions(t *testing.T) {
	client := Client{}
	actions := map[string]func(context.Context, string) error{
		"rectify": client.RectifyZone, "notify": client.NotifyZone, "axfr": client.AxfrZone,
	}
	for name, action := range actions {
		t.Run(name, func(t *testing.T) {
			if err := action(t.Context(), exampleZoneID); !errors.Is(err, connector.ErrUnsupported) {
				t.Fatalf("action error = %v, want unsupported", err)
			}
		})
	}
}
