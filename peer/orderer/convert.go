package orderer

import (
	"encoding/json"
	"madledger/core"
	pb "madledger/protos"
)

// This file convert structs and protos

// ConvertBlockFromTypesToPb convert block from core.Block to pb.Block
func ConvertBlockFromTypesToPb(typesBlock *core.Block) (*pb.Block, error) {
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

// ConvertBlockFromPbToTypes convert block from pb.Block to core.Block
func ConvertBlockFromPbToTypes(pbBlock *pb.Block) (*core.Block, error) {
	data, err := json.Marshal(pbBlock)
	if err != nil {
		return nil, err
	}
	var typesBlock core.Block
	err = json.Unmarshal(data, &typesBlock)
	if err != nil {
		return nil, err
	}
	return &typesBlock, nil
}

// ConvertTxFromPbToTypes convert Tx from pb.Tx to core.Tx
func ConvertTxFromPbToTypes(pbTx *pb.Tx) (*core.Tx, error) {
	data, err := json.Marshal(pbTx)
	if err != nil {
		return nil, err
	}
	var typesTx core.Tx
	err = json.Unmarshal(data, &typesTx)
	if err != nil {
		return nil, err
	}
	return &typesTx, nil
}

// ConvertTxFromTypesToPb convert Tx from core.Tx to pb.Tx
func ConvertTxFromTypesToPb(typesTx *core.Tx) (*pb.Tx, error) {
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
