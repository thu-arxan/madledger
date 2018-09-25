package protos

import (
	"encoding/json"
	"madledger/core/types"
)

// ConvertToTypes convert pb.Tx to types.Tx
func (tx *Tx) ConvertToTypes() (*types.Tx, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	var typesTx types.Tx
	err = json.Unmarshal(data, &typesTx)
	if err != nil {
		return nil, err
	}
	return &typesTx, nil
}

// NewTx is the constructor of Tx
func NewTx(typesTx *types.Tx) (*Tx, error) {
	data, err := json.Marshal(typesTx)
	if err != nil {
		return nil, err
	}
	var tx Tx
	err = json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
