package solo

import "madledger/consensus"

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
	return c.manager.start()
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return nil
}

// SyncBlocks is the implementation of interface
func (c *Consensus) SyncBlocks(channelID string, ch *chan consensus.Block) error {
	return nil
}

// GetNumber is the implementation of interface
func (c *Consensus) GetNumber(channelID string) (uint64, error) {
	return 0, nil
}
