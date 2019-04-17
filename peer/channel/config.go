package channel

import (
	"encoding/json"
	"errors"
	cc "madledger/blockchain/config"
	"madledger/core/types"
	"madledger/peer/db"
)

// AddConfigBlock add a config block
func (m *Manager) AddConfigBlock(block *types.Block) error {
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
		} else {
			status.Err = err.Error()
		}
		m.db.SetTxStatus(tx, status)
	}
	return nil
}

func getConfigPayload(tx *types.Tx) (*cc.Payload, error) {
	if tx.Data.ChannelID != types.CONFIGCHANNELID {
		return nil, errors.New("The tx does not belog to global channel")
	}
	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
