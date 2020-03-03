package channel

import (
	"madledger/core"
)

// AddGlobalBlock add a global block
func (m *Manager) AddGlobalBlock(block *core.Block) error {
	nums := make(map[string]uint64)
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		nums[payload.ChannelID] = payload.Num
	}
	m.coordinator.Unlocks(nums)
	wb := m.db.NewWriteBatch()
	wb.PutBlock(block)
	wb.Sync()
	return nil
}
