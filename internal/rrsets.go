package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

func (c *Client) GetRRset(ctx context.Context, zoneID, name, rrType string) (connector.RRset, error) {
	record := rrset()
	if zoneID != exampleZoneID || name != record.Name || rrType != record.Type {
		return connector.RRset{}, connector.ErrNotFound
	}
	return record, nil
}

func (c *Client) ReplaceRRset(ctx context.Context, zoneID string, input connector.RRsetReplaceInput) error {
	return connector.ErrUnsupported
}
func (c *Client) DeleteRRset(ctx context.Context, zoneID, name, rrType string) error {
	return connector.ErrUnsupported
}

func rrset() connector.RRset {
	return connector.RRset{
		Name: "www." + exampleZone, Type: "A", TTL: 300,
		Records: []connector.Record{{Content: "192.0.2.1"}},
	}
}
