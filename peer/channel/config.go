package channel

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/common/crypto"
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

				/* TODO: Gas
				gasprice = payload.gasprice
				ratio = payload.ratio
				maxgas = payload.maxgas
				这里应该将三个值都set到wb里去，需要新加db函数

				这里还应该进行token的分配，
				目前的一种想法是从sender账户里扣一个固定值，然后均分到每个member那里去。也需要两个db函数

				*** 问题：那么除了这个初始的token分配，之后的token分配该在哪里做？***
				*/
				maxgas := make([]byte, 8)
				binary.BigEndian.PutUint64(maxgas, payload.MaxGas)
				wb.Put(util.BytesCombine([]byte(channelID), []byte("maxgas")), maxgas)
				ratio := make([]byte, 8)
				binary.BigEndian.PutUint64(ratio, payload.AssetTokenRatio)
				wb.Put(util.BytesCombine([]byte(channelID), []byte("ratio")), ratio)
				gasprice := make([]byte, 8)
				binary.BigEndian.PutUint64(gasprice, payload.GasPrice)
				wb.Put(util.BytesCombine([]byte(channelID), []byte("gasprice")), gasprice)

				// sub 10000000 from sender's asset
				log.Debug("sub 10000000 from sender's asset")
				sender, _ := tx.GetSender()
				account, err := m.db.GetOrCreateAccount(sender)
				if err != nil {
					return err
				}
				if err = account.SubBalance(2000000000); err != nil {
					return err
				}
				wb.UpdateAccounts(account)

				log.Debug("give 10000000 to all the members and admin of this channel in token")
				// give 10000000 to all the members and admin of this channel in token
				part := uint64(2000000000/(len(payload.Profile.Admins)+len(payload.Profile.Members))) * payload.AssetTokenRatio
				var buf = make([]byte, 8)
				binary.BigEndian.PutUint64(buf, uint64(part))
				for _, admin := range payload.Profile.Admins {
					pk, _ := crypto.NewPublicKey(admin.PK)
					addr, _ := pk.Address()
					key := util.BytesCombine([]byte("token"), []byte(channelID), addr.Bytes())
					log.Debugf("config.go: the key is %v", key)
					wb.Put(key, buf)
				}
				for _, member := range payload.Profile.Members {
					// 现在sender 被减掉asset，key是address；member被加上token，key是token+channelID+addr
					pk, _ := crypto.NewPublicKey(member.PK)
					addr, _ := pk.Address()
					wb.Put(util.BytesCombine([]byte("token"), []byte(m.id), addr.Bytes()), buf)
				}

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
		} else {
			status.Err = err.Error()
		}
		log.Infof("config going to set status: channel: %s, tx: %s", m.id, tx.ID)
		wb.SetTxStatus(tx, status)
	}
	wb.PutBlock(block)
	wb.Sync()
	m.coordinator.Unlocks(nums)
	return nil
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
