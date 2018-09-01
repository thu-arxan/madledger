package server

import (
	"encoding/json"
	"madledger/core/types"
	pb "madledger/protos"
)

// This file convert structs and protos

// ConvertBlockFromTypesToPb convert block from types.Block to pb.Block
func ConvertBlockFromTypesToPb(typesBlock *types.Block) (*pb.Block, error) {
	bytes, err := json.Marshal(typesBlock)
	if err != nil {
		return nil, err
	}
	var pbBlock pb.Block
	err = json.Unmarshal(bytes, &pbBlock)
	if err != nil {
		return nil, err
	}
	return &pbBlock, nil
}

// ConvertBlockFromPbToTypes convert block from pb.Block to types.Block
func ConvertBlockFromPbToTypes(pbBlock *pb.Block) (*types.Block, error) {
	bytes, err := json.Marshal(pbBlock)
	if err != nil {
		return nil, err
	}
	var typesBlock types.Block
	err = json.Unmarshal(bytes, &typesBlock)
	if err != nil {
		return nil, err
	}
	return &typesBlock, nil
}
