package raft

import "encoding/json"

// Tx is the union of ChannelID and Data
type Tx struct {
	ChannelID string
	Data      []byte
}

// NewTx is the constructor of Tx
func NewTx(channelID string, data []byte) *Tx {
	return &Tx{
		ChannelID: channelID,
		Data:      data,
	}
}

// Bytes converty tx to bytes
func (tx *Tx) Bytes() []byte {
	data, _ := json.Marshal(tx)
	return data
}
