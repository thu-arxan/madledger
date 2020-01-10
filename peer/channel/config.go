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
	for i, tx := range block.Transactions {
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: block.Header.Number,
			BlockIndex:  i,
			Output:      nil,
		}
		payload, err := getConfigPayload(tx)
		if err == nil {
			channelID := payload.ChannelID
			if payload.Profile.Public {
				m.db.AddChannel(channelID)
			} else {
				var remove = true
				for _, member := range payload.Profile.Members {
					if member.Equal(m.identity) {
						m.db.AddChannel(channelID)
						remove = false
						break
					}
				}
				if remove && m.db.BelongChannel(channelID) {
					m.db.DeleteChannel(channelID)
				}
			}
			nums[payload.ChannelID] = 0
		} else {
			status.Err = err.Error()
		}
		m.db.SetTxStatus(tx, status)
	}
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
