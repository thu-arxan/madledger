package tendermint

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"

	node "github.com/tendermint/tendermint/node"
)

// Node obtains a tendermint node
type Node struct {
	// tendermintNode
	tn *node.Node
}

// NewNode is the constructor of Node
func NewNode(cfg *Config, app abci.Application) (*Node, error) {
	n := new(Node)
	logger := NewLogger()

	conf := config.DefaultConfig()
	conf.RootDir = cfg.Dir
	conf.Consensus.RootDir = cfg.Dir
	conf.Mempool.RootDir = cfg.Dir
	conf.P2P.RootDir = cfg.Dir
	conf.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.P2P)
	conf.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.RPC)
	conf.BaseConfig.ProxyApp = fmt.Sprintf("tcp://0.0.0.0:%d", cfg.Port.App)
	conf.Consensus.CreateEmptyBlocks = false
	conf.P2P.PersistentPeers = strings.Join(cfg.P2PAddress, ",")

	tn, err := node.DefaultNewNode(conf, logger)
	if err != nil {
		return n, err
	}
	n.tn = tn

	return n, nil
}

// Start runs the node
func (n *Node) Start() error {
	err := n.tn.Start()
	if err != nil {
		return err
	}
	n.tn.RunForever()
	return nil
}
