package channel

import (
	"madledger/core"
)

//todo: ab
// AddAssetBlock add an asset block
func (m *Manager) AddAssetBlock(block *core.Block) error {
	//nums := make(map[string]uint64)
	//for _, tx := range block.Transactions {
	//	payload, err := tx.GetAccountTxPayload()
	//	if err != nil {
	//		return err
	//	}
	//	nums[payload.ChannelID] = payload.Num
	//}
	//m.coordinator.Unlocks(nums)
	return nil
}
