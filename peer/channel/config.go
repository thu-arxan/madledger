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
	nums := make(map[string]uint64)
	wb := m.db.NewWriteBatch()
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
				nums[payload.ChannelID] = 0
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

func getConfigPayload(tx *core.Tx) (*cc.Payload, error) {
	if tx.Data.ChannelID != core.CONFIGCHANNELID {
		return nil, errors.New("The tx does not belog to global channel")
	}
	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
