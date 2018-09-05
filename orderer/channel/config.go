package channel

import (
	"encoding/json"
	cc "madledger/blockchain/config"
	"madledger/core/types"
)

// AddConfigBlock add a config block
// The block is formated, so there is no need to verify
func (manager *Manager) AddConfigBlock(block *types.Block) error {
	for _, tx := range block.Transactions {
		var payload cc.Payload
		json.Unmarshal(tx.Data.Payload, &payload)
		err := manager.db.UpdateChannel(payload.ChannelID, payload.Profile)
		if err != nil {
			return err
		}
	}
	return nil
}
