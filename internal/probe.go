package internal

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

// Probe demonstrates a bounded upstream request without reading any credentials.
func (c *Client) Probe(ctx context.Context) (connector.Profile, error) {
	response, err := c.exchange(ctx, abi.HTTPRequest{
		Method: http.MethodGet, Parts: []string{"health"},
		Headers: http.Header{"Accept": {"application/json"}}, Limit: 4096,
	})
	if err != nil {
		return connector.Profile{}, err
	}
	switch response.Status {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		return connector.Profile{}, connector.ErrAuthentication
	default:
		return connector.Profile{}, connector.ErrUnavailable
	}
	var health struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(response.Body, &health); err != nil || health.Version == "" {
		return connector.Profile{}, connector.ErrMalformed
	}
	return connector.Profile{DaemonType: providerName, Version: health.Version}, nil
}
