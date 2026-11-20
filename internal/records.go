package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

// This connector operates on RRsets. Individual-record operations have no
// upstream HTTP routes and report the unsupported capability.
func (c *Client) ListRecords(ctx context.Context, zoneID string) ([]connector.IndividualRecord, error) {
	return nil, connector.ErrUnsupported
}

func (c *Client) GetRecord(ctx context.Context, zoneID, recordID string) (connector.IndividualRecord, error) {
	return connector.IndividualRecord{}, connector.ErrUnsupported
}

func (c *Client) CreateRecord(ctx context.Context, zoneID string, input connector.IndividualRecordInput) (connector.IndividualRecord, error) {
	return connector.IndividualRecord{}, connector.ErrUnsupported
}

func (c *Client) UpdateRecord(ctx context.Context, zoneID, recordID string, input connector.IndividualRecordInput) (connector.IndividualRecord, error) {
	return connector.IndividualRecord{}, connector.ErrUnsupported
}

func (c *Client) DeleteRecord(ctx context.Context, zoneID, recordID string) error {
	return connector.ErrUnsupported
}

var _ connector.IndividualRecordOperations = (*Client)(nil)
