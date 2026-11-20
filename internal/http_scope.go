package internal

import (
	"context"
	"net/http"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

// scopeHTTP keeps each DNS operation within its upstream routes. Add a rule
// alongside every operation that starts making requests.
func scopeHTTP(operation abi.Operation, exchange Exchange) Exchange {
	if exchange == nil {
		return nil
	}
	return func(ctx context.Context, request abi.HTTPRequest) (abi.HTTPResponse, error) {
		if operation != abi.Probe || request.Method != http.MethodGet ||
			len(request.Parts) != 1 || request.Parts[0] != "health" {
			return abi.HTTPResponse{}, connector.ErrValidation
		}
		return exchange(ctx, request)
	}
}
