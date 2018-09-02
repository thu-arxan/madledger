package server

import (
	"encoding/json"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core/types"
	pb "madledger/protos"
)

// This file convert structs and protos

// ConvertBlockFromTypesToPb convert block from types.Block to pb.Block
func ConvertBlockFromTypesToPb(typesBlock *types.Block) (*pb.Block, error) {
	header, err := ConvertBlockHeaderFromTypesToPb(typesBlock.Header)
	if err != nil {
		return nil, err
	}
	var transactions []*pb.Tx
	for _, typesTx := range typesBlock.Transactions {
		pbTx, err := ConvertTxFromTypesToPb(typesTx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, pbTx)
	}
	return &pb.Block{
		Header:       header,
		Transactions: transactions,
	}, nil
}

// ConvertBlockFromPbToTypes convert block from pb.Block to types.Block
func ConvertBlockFromPbToTypes(pbBlock *pb.Block) (*types.Block, error) {
	header, err := ConvertBlockHeaderFromPbToTypes(pbBlock.Header)
	if err != nil {
		return nil, err
	}
	var transactions []*types.Tx
	for _, pbTx := range pbBlock.Transactions {
		typesTx, err := ConvertTxFromPbToTypes(pbTx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, typesTx)
	}
	return &types.Block{
		Header:       header,
		Transactions: transactions,
	}, nil
}

//ConvertBlockHeaderFromPbToTypes convert BlockHeader form pb.BlockHeader to types.BlockHeader
func ConvertBlockHeaderFromPbToTypes(pbBlockHeader *pb.BlockHeader) (*types.BlockHeader, error) {
	bytes, err := json.Marshal(pbBlockHeader)
	if err != nil {
		return nil, err
	}
	var typesBlockHeader types.BlockHeader
	err = json.Unmarshal(bytes, &typesBlockHeader)
	if err != nil {
		return nil, err
	}
	return &typesBlockHeader, nil
}

//ConvertBlockHeaderFromTypesToPb convert BlockHeader form types.BlockHeader to pb.BlockHeader
func ConvertBlockHeaderFromTypesToPb(typesBlockHeader *types.BlockHeader) (*pb.BlockHeader, error) {
	bytes, err := json.Marshal(typesBlockHeader)
	if err != nil {
		return nil, err
	}
	var pbBlockHeader pb.BlockHeader
	err = json.Unmarshal(bytes, &pbBlockHeader)
	if err != nil {
		return nil, err
	}
	return &pbBlockHeader, nil
}

// ConvertTxFromPbToTypes convert Tx from pb.Tx to types.Tx
func ConvertTxFromPbToTypes(tx *pb.Tx) (*types.Tx, error) {
	txData, err := ConvertTxDataFromPbToTytes(tx.Data)
	if err != nil {
		return nil, err
	}
	return &types.Tx{
		Data: txData,
		Time: tx.Time,
	}, nil
}

// ConvertTxFromTypesToPb convert Tx from types.Tx to pb.Tx
func ConvertTxFromTypesToPb(tx *types.Tx) (*pb.Tx, error) {
	txData, err := ConvertTxDataFromTypesToPb(tx.Data)
	if err != nil {
		return nil, err
	}
	return &pb.Tx{
		Data: txData,
		Time: tx.Time,
	}, nil
}

// ConvertTxDataFromPbToTytes convert TxData from pb.TxData to types.TxData
func ConvertTxDataFromPbToTytes(pbTxData *pb.TxData) (*types.TxData, error) {
	recipient, err := common.AddressFromBytes(pbTxData.Recipient)
	if err != nil {
		return nil, err
	}
	var typesTxData = types.TxData{
		ChannelID:    pbTxData.ChannelID,
		AccountNonce: pbTxData.AccountNonce,
		Recipient:    recipient,
		Payload:      pbTxData.Payload,
		Version:      pbTxData.Version,
	}
	if pbTxData.Sig == nil {
		return &typesTxData, nil
	}
	pk, err := crypto.NewPublicKey(pbTxData.Sig.PK)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.NewSignature(pbTxData.Sig.Sig)
	if err != nil {
		return nil, err
	}
	typesTxData.Sig = &types.TxSig{
		PK:  pk,
		Sig: sig,
	}

	return &typesTxData, nil
}

// ConvertTxDataFromTypesToPb convert TxData from types.TxData to pb.TxData
func ConvertTxDataFromTypesToPb(typesTxData *types.TxData) (*pb.TxData, error) {
	recipient := typesTxData.Recipient.Bytes()
	var pbTxData = pb.TxData{
		ChannelID:    typesTxData.ChannelID,
		AccountNonce: typesTxData.AccountNonce,
		Recipient:    recipient,
		Payload:      typesTxData.Payload,
		Version:      typesTxData.Version,
	}
	if typesTxData.Sig == nil {
		return &pbTxData, nil
	}
	pk, err := typesTxData.Sig.PK.Bytes()
	if err != nil {
		return nil, err
	}
	sig, err := typesTxData.Sig.Sig.Bytes()
	if err != nil {
		return nil, err
	}
	pbTxData.Sig = &pb.TxSig{
		PK:  pk,
		Sig: sig,
	}
	return &pbTxData, nil
}
