package protos

import (
	"encoding/json"
	"madledger/core/types"
)

// NewBlock is the constructor of Block
func NewBlock(typesBlock *types.Block) (*Block, error) {
	data, err := json.Marshal(typesBlock)
	if err != nil {
		return nil, err
	}
	var block Block
	err = json.Unmarshal(data, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// ConvertToTypes convert pb.Block to types.Block
func (block *Block) ConvertToTypes() (*types.Block, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	var typesBlock types.Block
	err = json.Unmarshal(data, &typesBlock)
	if err != nil {
		return nil, err
	}
	return &typesBlock, nil
}
