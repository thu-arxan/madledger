package tendermint

import (
	"errors"
	"fmt"
	"madledger/consensus"
	"madledger/core"
	"sync"

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
	lock   sync.Mutex
	status consensus.Status
	app    *Glue
	node   *Node
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
		status: consensus.Stopped,
		app:    app,
		node:   node,
	}, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Info("Trying to start consensus")
	err := c.app.Start()
	if err != nil {
		return err
	}
	err = c.node.Start()
	if err != nil {
		return err
	}
	log.Info("Start consensus...")
	c.status = consensus.Started

	return nil
}

// AddChannel add a channel
// Because we are not using multi-group to improve performance, so we can just ignore this function in tendermint
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	return nil
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(tx *core.Tx) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.status != consensus.Started {
		return errors.New("The service is not started")
	}

	bytes, _ := tx.Bytes()
	return c.app.AddTx(tx.Data.ChannelID, bytes)
}

// SyncBlocks is the implementation of interface
func (c *Consensus) SyncBlocks(channelID string, ch *chan consensus.Block) error {
	c.app.SetSyncChan(channelID, ch)
	return nil
}

// Stop is the implementation of interface
// todo: we need to make sure that the consensus will not provide service
func (c *Consensus) Stop() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.status = consensus.Stopped
	c.node.Stop()
	c.app.Stop()
	return nil
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return c.app.GetBlock(channelID, num, async)
}
