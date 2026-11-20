package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

const (
	exampleZoneID = "example-zone"
	exampleZone   = "example.test."
)

func (c *Client) ListZones(ctx context.Context) ([]connector.Zone, error) {
	return []connector.Zone{zone()}, nil
}

func (c *Client) GetZone(ctx context.Context, id string) (connector.Zone, error) {
	if id != exampleZoneID {
		return connector.Zone{}, connector.ErrNotFound
	}
	return zone(), nil
}

func (c *Client) GetZoneMetadata(ctx context.Context, id string) (connector.ZoneMetadata, error) {
	if id != exampleZoneID {
		return connector.ZoneMetadata{}, connector.ErrNotFound
	}
	return connector.ZoneMetadata{ID: exampleZoneID, Name: exampleZone, Kind: "example"}, nil
}

func zone() connector.Zone {
	return connector.Zone{
		ID: exampleZoneID, Name: exampleZone, Kind: "example", RRsets: []connector.RRset{rrset()},
		Capabilities: &connector.ZoneCapabilities{
			RecordTypes: []string{"A"}, DNSSEC: unsupportedDNSSEC(),
		},
	}
}
