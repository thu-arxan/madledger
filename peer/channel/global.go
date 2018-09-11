package channel

import "madledger/core/types"

// AddGlobalBlock add a global block
func (m *Manager) AddGlobalBlock(block *types.Block) error {
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		log.Infof("Channel %s, block %d", payload.ChannelID, payload.Num)
	}
	return nil
}
