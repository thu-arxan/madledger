package tendermint

import (
	"fmt"

	rc "github.com/tendermint/tendermint/rpc/client"
)

// Client contain the rpcclient of tendermint
type Client struct {
	tc *rc.HTTP
}

// NewClient is the constructor of Client
func NewClient(port int) (*Client, error) {
	client := rc.NewHTTP(fmt.Sprintf("tcp://0.0.0.0:%d", port), "/websocket")

	return &Client{
		tc: client,
	}, nil
}

// AddTx send a tx into tendermint network
func (c *Client) AddTx(tx []byte) error {
	_, err := c.tc.BroadcastTxSync(tx)
	return err
}
