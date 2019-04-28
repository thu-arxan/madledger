package tendermint

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	tc "github.com/tendermint/tendermint/config"

	node "github.com/tendermint/tendermint/node"
)

// Node obtains a tendermint node
type Node struct {
	// tendermintNode
	tn   *node.Node
	conf *tc.Config
}

// NewNode is the constructor of Node
func NewNode(cfg *Config, app abci.Application) (*Node, error) {
	n := new(Node)

	conf := tc.DefaultConfig()
	conf.SetRoot(cfg.Dir)
	conf.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.P2P)
	conf.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.RPC)
	conf.BaseConfig.ProxyApp = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.App)
	conf.Consensus.CreateEmptyBlocks = false
	conf.P2P.AddrBookStrict = false
	conf.P2P.AllowDuplicateIP = true
	conf.P2P.PersistentPeers = strings.Join(cfg.P2PAddress, ",")

	n.conf = conf
	return n, nil
}

// Start runs the node
func (n *Node) Start() error {
	logger := NewLogger()
	tn, err := node.DefaultNewNode(n.conf, logger)
	if err != nil {
		return err
	}
	n.tn = tn
	err = n.tn.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the node
func (n *Node) Stop() {
	if n.tn != nil {
		n.tn.Stop()
	}
}
