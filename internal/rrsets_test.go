package internal

import (
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/connector"
)

func TestGetRRsetMatchesNameAndType(t *testing.T) {
	client := Client{}
	name := "www." + exampleZone
	if _, err := client.GetRRset(t.Context(), exampleZoneID, name, "AAAA"); !errors.Is(err, connector.ErrNotFound) {
		t.Fatalf("GetRRset() with different type = %v, want not found", err)
	}
	records, err := client.GetRRset(t.Context(), exampleZoneID, name, "A")
	if err != nil || len(records.Records) != 1 || records.Records[0].Content != "192.0.2.1" {
		t.Fatalf("GetRRset() = %+v, %v", records, err)
	}
}

func TestReadOnlyConnectorRejectsRRsetWrites(t *testing.T) {
	client := Client{}
	input := connector.RRsetReplaceInput{Name: "www.example.test.", Type: "A", TTL: 300}
	if err := client.ReplaceRRset(t.Context(), exampleZoneID, input); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("ReplaceRRset() = %v, want unsupported", err)
	}
	if err := client.DeleteRRset(t.Context(), exampleZoneID, input.Name, input.Type); !errors.Is(err, connector.ErrUnsupported) {
		t.Fatalf("DeleteRRset() = %v, want unsupported", err)
	}
}
