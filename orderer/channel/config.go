package channel

import (
	"encoding/json"
	cc "madledger/blockchain/config"
	"madledger/consensus"
	"madledger/core/types"
)

// AddConfigBlock add a config block
// The block is formated, so there is no need to verify
func (manager *Manager) AddConfigBlock(block *types.Block) error {
	if block.Header.Number == 0 {
		return nil
	}
	for _, tx := range block.Transactions {
		var payload cc.Payload
		json.Unmarshal(tx.Data.Payload, &payload)
		var channelID = payload.ChannelID
		// This is a create channel tx
		if !manager.db.HasChannel(channelID) {
			// then start the consensus
			err := manager.coordinator.Consensus.AddChannel(channelID, consensus.DefaultConfig())
			channel, err := NewManager(channelID, manager.coordinator)
			if err != nil {
				return err
			}
			// create genesis block here
			// The genesis only contain the create tx now.
			genesisBlock := types.NewBlock(channelID, 0, types.GenesisBlockPrevHash, []*types.Tx{tx})
			err = channel.AddBlock(genesisBlock)
			if err != nil {
				return err
			}
			// then start the channel
			go func() {
				channel.Start()
			}()
			manager.coordinator.Managers[channelID] = channel
		}
		err := manager.db.UpdateChannel(channelID, payload.Profile)
		if err != nil {
			return err
		}
	}
	return nil
}
