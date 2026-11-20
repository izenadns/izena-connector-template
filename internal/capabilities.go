package internal

import (
	"context"

	"github.com/izenadns/izena-connector-contract/connector"
)

func (c *Client) Capabilities(ctx context.Context) (connector.ServerCapabilities, error) {
	return connector.ServerCapabilities{
		Provider: providerID, RecordModel: connector.RecordModelRRsets, RecordTypes: []string{"A"}, DNSSEC: unsupportedDNSSEC(),
	}, nil
}

func unsupportedDNSSEC() connector.DNSSECCapabilities {
	return connector.DNSSECCapabilities{Signing: connector.SigningUnsupported}
}
