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
	"encoding/binary"
	"encoding/json"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/common/util"
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
		if err == nil {
			switch len(payload.ChannelID) {
			case 0:
				log.Warnf("Fatal error! Nil channel id in config block, num: %d, index: %d", block.GetNumber(), i)
			default:
				channelID := payload.ChannelID

				// 创建通道时指定关于gas的设置，目前是gas price, ratio, gas limit
				if tx.GetReceiver().String() == core.CreateChannelContractAddress.String() {
					setChannelConfig(channelID, wb, payload)

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
				}

				nums[payload.ChannelID] = []uint64{0}
			}
		} else {
			status.Err = err.Error()
		}
		wb.SetTxStatus(tx, status)
	}
	wb.PutBlock(block)
	wb.Sync()
	m.coordinator.Unlocks(nums)
	return nil
}

// keys are unified now in orderer / peer
// if needed to modify, modify all of them
func setChannelConfig(channelID string, wb db.WriteBatch, payload *cc.Payload) {
	ratioBytes := make([]byte, 8)
	ratio := payload.AssetTokenRatio
	if ratio == 0 {
		ratio = 1
	}
	binary.BigEndian.PutUint64(ratioBytes, ratio)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("ratio")), ratioBytes)

	// gasPrice could be zero
	gasPriceBytes := make([]byte, 8)
	gasPrice := payload.GasPrice
	binary.BigEndian.PutUint64(gasPriceBytes, gasPrice)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("gasPrice")), gasPriceBytes)

	maxGasBytes := make([]byte, 8)
	maxGas := payload.MaxGas
	if maxGas == 0 {
		maxGas = 1000000
	}
	binary.BigEndian.PutUint64(maxGasBytes, maxGas)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("maxGas")), maxGasBytes)

}

func getConfigPayload(tx *core.Tx) (*cc.Payload, error) {
	if tx.Data.ChannelID != core.CONFIGCHANNELID {
		return nil, errors.New("The tx does not belong to global channel")
	}
	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
