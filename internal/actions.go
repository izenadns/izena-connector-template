package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

func (c *Client) RectifyZone(ctx context.Context, zoneID string) error {
	return connector.ErrUnsupported
}
func (c *Client) NotifyZone(ctx context.Context, zoneID string) error {
	return connector.ErrUnsupported
}
func (c *Client) AxfrZone(ctx context.Context, zoneID string) error {
	return connector.ErrUnsupported
}
