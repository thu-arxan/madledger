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
// TODO: We should read the document to understand the return of BroadcastTxSync
// Note: BroadcastTxSync may return error because the tx in the cache aleardy
func (c *Client) AddTx(tx []byte) error {
	_, err := c.tc.BroadcastTxSync(tx)
	//broadcast_tx_sync: Response error: RPC error -32603 - Internal error: EOF
	if err != nil && err.Error() != "broadcast_tx_sync: Response error: RPC error -32603 - Internal error: Tx already exists in cache" {
		log.Infof("AddTx meets an error: %s\n", err)
	}

	return nil
}
