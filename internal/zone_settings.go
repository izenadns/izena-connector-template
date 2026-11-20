package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

func (c *Client) GetZoneSettings(ctx context.Context, zoneID string) (connector.ZoneSettings, error) {
	return connector.ZoneSettings{}, connector.ErrUnsupported
}

func (c *Client) PatchZoneSettings(ctx context.Context, zoneID string, patch connector.ZoneSettingsPatch) error {
	return connector.ErrUnsupported
}

var _ connector.ZoneSettingsOperations = (*Client)(nil)
