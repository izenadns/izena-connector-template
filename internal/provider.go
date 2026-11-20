// Package internal implements the provider contract and upstream operations.
package internal

import (
	"context"
	"net/url"
	"strings"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

const (
	providerID       = "example-dns"
	providerName     = "Example DNS"
	baseURLKey       = "endpoint"
	credentialKey    = "token"
	maxBaseURLLength = 2048
)

// Exchange delegates a request to Izena's connection-scoped HTTP host.
type Exchange func(context.Context, abi.HTTPRequest) (abi.HTTPResponse, error)

type Provider struct {
	HTTP Exchange
}

func (*Provider) Describe() connector.Description {
	return connector.Description{
		ID: providerID, Name: providerName,
		ConnectionFields: []connector.ConnectionField{
			{Key: baseURLKey, Label: "API base URL", URL: true, Required: true, MaxLength: maxBaseURLLength,
				Placeholder: "https://dns.example.test"},
			{Key: credentialKey, Label: "Bearer token", Secret: true, Editable: true, MaxLength: 4096},
		},
	}
}

func (p *Provider) Manifest() abi.Manifest {
	return abi.Manifest{
		Description: p.Describe(),
		HTTP: abi.HTTPBinding{
			URLSetting:  baseURLKey,
			Credentials: []abi.HTTPHeaderCredential{{Header: "Authorization", SecretField: credentialKey, Prefix: "Bearer "}},
		},
	}
}

func (*Provider) PrepareSettings(_ context.Context, settings map[string]string) (connector.PreparedSettings, error) {
	endpoint, err := url.Parse(strings.TrimSpace(settings[baseURLKey]))
	if err != nil || endpoint == nil || endpoint.Hostname() == "" || endpoint.User != nil || endpoint.Opaque != "" ||
		(endpoint.Scheme != "http" && endpoint.Scheme != "https") ||
		endpoint.RawQuery != "" || endpoint.ForceQuery || endpoint.Fragment != "" {
		return connector.PreparedSettings{}, connector.ErrInvalidServerURL
	}
	baseURL := strings.TrimRight(endpoint.String(), "/")
	normalized := map[string]string{baseURLKey: baseURL}
	return connector.PreparedSettings{Settings: normalized, TargetKey: normalized[baseURLKey]}, nil
}

func (p *Provider) Connect(config connector.ConnectionConfig) (connector.Client, error) {
	client, err := p.connect(config.Settings)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (p *Provider) connect(settings map[string]string) (*Client, error) {
	if p.HTTP == nil {
		return nil, connector.ErrUnavailable
	}
	if settings[baseURLKey] == "" {
		return nil, connector.ErrInvalidServerURL
	}
	return &Client{exchange: p.HTTP}, nil
}

var _ connector.Provider = (*Provider)(nil)
