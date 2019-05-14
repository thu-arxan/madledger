package tendermint

import (
	"fmt"
	"os"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	tc "github.com/tendermint/tendermint/config"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"

	node "github.com/tendermint/tendermint/node"
)

// Node obtains a tendermint node
type Node struct {
	// tendermintNode
	tn    *node.Node
	tnDBs map[string]dbm.DB
	conf  *tc.Config
	app   abci.Application
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
	// conf.Mempool.Broadcast = false
	// conf.P2P.Seeds = conf.P2P.PersistentPeers

	n.conf = conf

	n.app = app
	return n, nil
}

// Start runs the node
func (n *Node) Start() error {
	err := n.StartTendermintNode()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the node
func (n *Node) Stop() {
	if n.tn != nil {
		n.tn.Stop()
		n.tn.Wait()
		for _, db := range n.tnDBs {
			db.Close()
		}
	}
}

// StartTendermintNode is a copy of tendermint.DefaultNewNode
// The reason why we need this is that we want to close db connection
func (n *Node) StartTendermintNode() error {
	logger := NewLogger()
	config := n.conf

	n.tnDBs = make(map[string]dbm.DB)

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return err
	}

	oldPrivVal := config.OldPrivValidatorFile()
	newPrivValKey := config.PrivValidatorKeyFile()
	newPrivValState := config.PrivValidatorStateFile()
	if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
		oldPV, err := privval.LoadOldFilePV(oldPrivVal)
		if err != nil {
			return fmt.Errorf("Error reading OldPrivValidator from %s: %s", oldPrivVal, err)
		}
		oldPV.Upgrade(newPrivValKey, newPrivValState)
	}

	tn, err := node.NewNode(config,
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		// proxy.DefaultClientCreator(config.ProxyApp, config.ABCI, config.DBDir()),
		proxy.NewLocalClientCreator(n.app),
		node.DefaultGenesisDocProviderFunc(config),
		NewDBProvide(n),
		node.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)
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

// NewDBProvide will set node some dbs
func NewDBProvide(n *Node) func(ctx *node.DBContext) (dbm.DB, error) {
	return func(ctx *node.DBContext) (dbm.DB, error) {
		dbType := dbm.DBBackendType(ctx.Config.DBBackend)
		db := dbm.NewDB(ctx.ID, dbType, ctx.Config.DBDir())
		n.tnDBs[ctx.ID] = db
		return db, nil
	}
}
