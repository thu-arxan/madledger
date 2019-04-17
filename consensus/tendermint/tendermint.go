package tendermint

import (
	"fmt"
	"madledger/consensus"
	"time"

	"github.com/sirupsen/logrus"
)

// Note: We just reach consensus in one chain of tendermint now.
// So we will not distinguish channel in the consensus layer and we consensus on
// one block then we split it into different blocks of different channels.

var (
	log = logrus.WithFields(logrus.Fields{"app": "consensus", "package": "tendermint"})
)

// Consensus is the implementation of tendermint
type Consensus struct {
	app  *Glue
	node *Node
}

// NewConsensus is the constructor of tendermint.Consensus
// Note: We are not going to support different configs of different channels.
// TODO: Not finished yet
func NewConsensus(channels map[string]consensus.Config, cfg *Config) (consensus.Consensus, error) {
	app, err := NewGlue(fmt.Sprintf("%s/.glue", cfg.Dir), &cfg.Port)
	if err != nil {
		return nil, err
	}
	node, err := NewNode(cfg, app)
	if err != nil {
		return nil, err
	}
	return &Consensus{
		app:  app,
		node: node,
	}, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	log.Info("Trying to start consensus")
	go c.app.Start()
	time.Sleep(200 * time.Millisecond)
	go c.node.Start()
	time.Sleep(300 * time.Millisecond)
	log.Info("Start consensus...")
	return nil
}

// AddChannel add a channel
// Because we are not using multi-group to improve performance, so we can just ignore this function in tendermint
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	return nil
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return c.app.AddTx(channelID, tx)
}

// SyncBlocks is the implementation of interface
func (c *Consensus) SyncBlocks(channelID string, ch *chan consensus.Block) error {
	c.app.SetSyncChan(channelID, ch)
	return nil
}

// Stop is the implementation of interface
// todo: implement the stop function
func (c *Consensus) Stop() error {
	return nil
}

// GetBlock is the implementation of interface
// TODO: Implementation it.
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return c.app.GetBlock(channelID, num, async)
}
