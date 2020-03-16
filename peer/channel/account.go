// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"madledger/core"
)

// AddAssetBlock add an account block
// todo: ab
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
