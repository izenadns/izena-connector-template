package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

func (c *Client) CreateZone(ctx context.Context, input connector.ZoneCreateInput) (connector.Zone, error) {
	return connector.Zone{}, connector.ErrUnsupported
}
func (c *Client) UpdateZone(ctx context.Context, id string, input connector.ZoneUpdateInput) (connector.ZoneUpdateResult, error) {
	return connector.ZoneUpdateResult{}, connector.ErrUnsupported
}
func (c *Client) DeleteZone(ctx context.Context, id string) error { return connector.ErrUnsupported }
