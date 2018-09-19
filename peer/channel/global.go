package channel

import (
	"madledger/core/types"
)

// AddGlobalBlock add a global block
func (m *Manager) AddGlobalBlock(block *types.Block) error {
	nums := make(map[string]uint64)
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		nums[payload.ChannelID] = payload.Num
	}
	m.coordinator.Unlocks(nums)
	return nil
}
