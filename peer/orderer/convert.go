package orderer

import (
	"encoding/json"
	"madledger/core/types"
	pb "madledger/protos"
)

// This file convert structs and protos

// ConvertBlockFromTypesToPb convert block from types.Block to pb.Block
func ConvertBlockFromTypesToPb(typesBlock *types.Block) (*pb.Block, error) {
	data, err := json.Marshal(typesBlock)
	if err != nil {
		return nil, err
	}
	var pbBlock pb.Block
	err = json.Unmarshal(data, &pbBlock)
	if err != nil {
		return nil, err
	}
	return &pbBlock, nil
}

// ConvertBlockFromPbToTypes convert block from pb.Block to types.Block
func ConvertBlockFromPbToTypes(pbBlock *pb.Block) (*types.Block, error) {
	data, err := json.Marshal(pbBlock)
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

// ConvertTxFromPbToTypes convert Tx from pb.Tx to types.Tx
func ConvertTxFromPbToTypes(pbTx *pb.Tx) (*types.Tx, error) {
	data, err := json.Marshal(pbTx)
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

// ConvertTxFromTypesToPb convert Tx from types.Tx to pb.Tx
func ConvertTxFromTypesToPb(typesTx *types.Tx) (*pb.Tx, error) {
	data, err := json.Marshal(typesTx)
	if err != nil {
		return nil, err
	}
	var pbTx pb.Tx
	err = json.Unmarshal(data, &pbTx)
	if err != nil {
		return nil, err
	}
	return &pbTx, nil
}
