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

// AddGlobalBlock add a global block
func (m *Manager) AddGlobalBlock(block *core.Block) error {
	nums := make(map[string][]uint64)
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		switch payload.ChannelID {
		case core.CONFIGCHANNELID, core.ASSETCHANNELID:
			m.coordinator.Unlocks(map[string][]uint64{payload.ChannelID: []uint64{payload.Num}})
			// zhq todo: am i doing it correct?
		default:
			nums[payload.ChannelID] = append(nums[payload.ChannelID], payload.Num)
		}
	}
	m.coordinator.Unlocks(nums)
	wb := m.db.NewWriteBatch()
	wb.PutBlock(block)
	wb.Sync()
	return nil
}
