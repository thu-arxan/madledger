package solo

import (
	"madledger/consensus"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "consensus", "package": "solo"})
)

// Consensus is the implementaion of solo consensus
type Consensus struct {
	manager *manager
}

// NewConsensus is the constructor of solo.Consensus
func NewConsensus(channels map[string]consensus.Config) (consensus.Consensus, error) {
	c := new(Consensus)
	c.manager = newManager()
	for id, cfg := range channels {
		err := c.manager.add(id, cfg)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	err := c.manager.start()
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

// AddChannel add a channel
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	err := c.manager.add(channelID, cfg)
	if err != nil {
		return err
	}
	return c.manager.startChannel(channelID)
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return c.manager.AddTx(channelID, tx)
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	channel, err := c.manager.get(channelID)
	if err != nil {
		return nil, err
	}
	return channel.getBlock(num, async)
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	return c.manager.stop()
}
