package solo

import (
	"madledger/consensus"
	"time"

	"github.com/rs/zerolog/log"
)

type channel struct {
	id     string
	config consensus.Config
	txChan chan []byte
	txs    [][]byte
	num    uint64
	// store all blocks, maybe gc is needed
	// todo
	blocks map[uint64]*Block
	notify *chan consensus.Block
}

func newChannel(id string, config consensus.Config, notify *chan consensus.Block) *channel {
	return &channel{
		id:     id,
		config: config,
		num:    config.Number,
		txChan: make(chan []byte),
		notify: notify,
		blocks: make(map[uint64]*Block),
	}
}

func (c *channel) start() {
	ticker := time.NewTicker(time.Duration(c.config.Timeout) * time.Millisecond)
	defer ticker.Stop()
	log.Info().Msgf("Channel %s start", c.id)
	for {
		select {
		case <-ticker.C:
			log.Info().Msgf("Channel %s tick", c.id)
			c.generateBlock()
		case tx := <-c.txChan:
			c.txs = append(c.txs, tx)
			if len(c.txs) >= c.config.MaxSize {
				c.generateBlock()
			}
		}
	}
}

func (c *channel) generateBlock() {
	if len(c.txs) == 0 {
		return
	}
	var txs [][]byte
	block := Block{
		channelID: c.id,
		num:       c.num,
		txs:       txs,
	}
	for _, tx := range c.txs {
		block.txs = append(block.txs, tx)
	}
	if c.notify != nil {
		*c.notify <- block
	}
	c.blocks[block.num] = &block
	c.txs = make([][]byte, 0)
	c.num++
}
