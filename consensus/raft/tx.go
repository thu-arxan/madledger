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

// UnmarshalTx unmarshal tx from bytes
func UnmarshalTx(bytes []byte) (*Tx, error) {
	var tx Tx
	err := json.Unmarshal(bytes, &tx)

	return &tx, err
}
