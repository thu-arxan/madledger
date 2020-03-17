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
				sender, _ := tx.GetSender()

				// 更新通道的关于gas的设置，目前是gas price, ratio, gas limit
				// TODO: 应该判断是否是通道管理员，应该用payload.IsAdmin来判断，但是不知道怎么生成传入的参数，所以暂时用cake替代
				if tx.GetReceiver().String() == core.CreateChannelContractAddress.String() {
					updateChannelConfig(channelID, wb, payload)

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

				// 如果AssetToDistribute被设置这里还应该进行token的分配，
				// 目前的做法是从sender账户里扣AssetToDistribute这么多的asset，然后换算成token均分到每个member那里去
				cake := tx.Data.Value
				// 如果没有设置这个变量，则Umarshal之后是0
				if cake > 0 {
					// 现在sender 被减掉asset，key是address；member被加上token，key是token+channelID+addr
					account, err := m.db.GetOrCreateAccount(sender)
					if err != nil {
						status.Err = err.Error()
						wb.SetTxStatus(tx, status)
						continue
					}
					if err = account.SubBalance(cake); err != nil {
						status.Err = err.Error()
						wb.SetTxStatus(tx, status)
						continue
					}
					wb.UpdateAccounts(account)
					peopleNum := uint64(len(payload.Profile.Admins) + len(payload.Profile.Members))
					ratioKey := util.BytesCombine([]byte(channelID), []byte("ratio"))
					ratioByte, err := m.db.Get(ratioKey)
					if err != nil {
						status.Err = err.Error()
						wb.SetTxStatus(tx, status)
						continue
					}
					ratio := uint64(binary.BigEndian.Uint64(ratioByte))
					part := (cake / peopleNum) * ratio
					var value = make([]byte, 8)
					binary.BigEndian.PutUint64(value, part)
					members := append(payload.Profile.Admins, payload.Profile.Members...)
					for _, member := range members {
						wb.Put(util.BytesCombine([]byte("token"), []byte(channelID), member.PK), value)
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

func updateChannelConfig(channelID string, wb db.WriteBatch, payload *cc.Payload) {
	maxgas := make([]byte, 8)
	binary.BigEndian.PutUint64(maxgas, payload.MaxGas)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("maxgas")), maxgas)
	ratio := make([]byte, 8)
	binary.BigEndian.PutUint64(ratio, payload.AssetTokenRatio)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("ratio")), ratio)
	gasprice := make([]byte, 8)
	binary.BigEndian.PutUint64(gasprice, payload.GasPrice)
	wb.Put(util.BytesCombine([]byte(channelID), []byte("gasprice")), gasprice)
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
