package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

// Advertise DNSSEC support in capabilities when implementing these operations.
func (c *Client) GetDNSSEC(ctx context.Context, zoneID string) (connector.DNSSEC, error) {
	return connector.DNSSEC{}, connector.ErrUnsupported
}
func (c *Client) SetDNSSEC(ctx context.Context, zoneID string, enabled bool) error {
	return connector.ErrUnsupported
}
