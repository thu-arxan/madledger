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
	"encoding/json"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/core"
	"madledger/peer/db"
)

// AddConfigBlock add a config block
func (m *Manager) AddConfigBlock(block *core.Block) error {
	wb := m.db.NewWriteBatch()
	nums := make(map[string][]uint64)
	for i, tx := range block.Transactions {
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: block.Header.Number,
			BlockIndex:  i,
			Output:      nil,
		}
		payload, err := getConfigPayload(tx)
		if err != nil {
			status.Err = err.Error()
			wb.SetTxStatus(tx, status)
			continue
		}
		if len(payload.ChannelID == 0) {
			log.Warnf("Fatal error! Nil channel id in config block, num: %d, index: %d", block.GetNumber(), i)
			continue
		}

		channelID := payload.ChannelID
		if tx.GetReceiver().String() == core.CreateChannelContractAddress.String() {

			if payload.Profile.Public {
				wb.AddChannel(channelID)
				m.coordinator.hub.Broadcast("update", Update{
					ID:     channelID,
					Remove: false,
				})
			} else {
				var remove = true
				for _, member := range payload.Profile.Members {
					if member.Equal(m.identity) {
						wb.AddChannel(channelID)
						m.coordinator.hub.Broadcast("update", Update{
							ID:     channelID,
							Remove: false,
						})
						remove = false
						break
					}
				}
				if remove && m.db.BelongChannel(channelID) {
					wb.DeleteChannel(channelID)
					m.coordinator.hub.Broadcast("update", Update{
						ID:     channelID,
						Remove: true,
					})
				}
			}
			nums[payload.ChannelID] = []uint64{0}
		}
		// todo:
		// in orderer this part does not use write batch
		err := m.db.UpdateChannel(channelID, payload.Profile)
		if err != nil {
			status.Err = err
		}
		wb.SetTxStatus(tx, status)
	}
	wb.PutBlock(block)
	wb.Sync()
	m.coordinator.Unlocks(nums)
	return nil
}

func getConfigPayload(tx *core.Tx) (*cc.Payload, error) {
	if tx.Data.ChannelID != core.CONFIGCHANNELID {
		return nil, errors.New("The tx does not belong to config channel")
	}
	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
