package internal

import (
	"github.com/izenadns/izena-connector-contract/connector"
)

type Client struct {
	exchange Exchange
}

var _ connector.Client = (*Client)(nil)
