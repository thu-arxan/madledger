// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tendermint

import (
	"errors"
	"fmt"

	rc "github.com/tendermint/tendermint/rpc/client"
)

// Client contain the rpcclient of tendermint
type Client struct {
	port int
	tc   *rc.HTTP
}

// NewClient is the constructor of Client
func NewClient(port int) (*Client, error) {
	client := rc.NewHTTP(fmt.Sprintf("tcp://0.0.0.0:%d", port), "/websocket")

	return &Client{
		port: port,
		tc:   client,
	}, nil
}

// AddTx send a tx into tendermint network
// TODO: We should read the document to understand the return of BroadcastTxSync
// Note: BroadcastTxSync may return error because the tx in the cache aleardy
func (c *Client) AddTx(tx []byte) error {
	//t, _ := BytesToTx(tx)
	//log.Infof("[%d]Here client add tx %s", c.port, string(t.Data))
	defer func() {
		//log.Infof("[%d]Done client add tx %s", c.port, string(t.Data))
	}()
	_, err := c.tc.BroadcastTxSync(tx)
	//broadcast_tx_sync: Response error: RPC error -32603 - Internal error: EOF
	if err != nil && err.Error() != "broadcast_tx_sync: Response error: RPC error -32603 - Internal error: Tx already exists in cache" {
		log.Infof("AddTx meets an error: %s\n", err)
		return errors.New("Meet rpc error")
	} else if err != nil {
		log.Infof("AddTx meets an error:%s", err.Error())
	}

	return nil
}
