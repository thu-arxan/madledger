package server

import (
	"errors"
	"fmt"
	"madledger/common"
	"madledger/core/types"
	pb "madledger/protos"
	"madledger/common/crypto"

	"golang.org/x/net/context"
)

// FetchBlock is the implementation of protos
func (s *Server) FetchBlock(ctx context.Context, req *pb.FetchBlockRequest) (*pb.Block, error) {
	block, err := s.cc.FetchBlock(req.ChannelID, req.Number, req.Behavior == pb.Behavior_RETURN_UNTIL_READY)
	if err != nil {
		return nil, err
	}
	return pb.NewBlock(block)
}

// ListChannels is the implementation of protos
func (s *Server) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	return s.cc.ListChannels(req)
}

// CreateChannel is the implementation of protos
func (s *Server) CreateChannel(ctx context.Context, req *pb.CreateChannelRequest) (*pb.ChannelInfo, error) {
	tx, err := req.GetTx().ConvertToTypes()
	if err != nil {
		return nil, err
	}
	if !tx.Verify() {
		return nil, errors.New("The tx is not a valid tx")
	}
	if tx.GetReceiver().String() != types.CreateChannelContractAddress.String() {
		return nil, errors.New("The receiver of the tx is not the valid contract address")
	}
	_, err = s.cc.CreateChannel(tx)
	if err != nil {
		return nil, err
	}
	return &pb.ChannelInfo{}, nil
}

// AddTx is the implementation of protos
func (s *Server) AddTx(ctx context.Context, req *pb.AddTxRequest) (*pb.TxStatus, error) {
	var status pb.TxStatus
	tx, err := req.Tx.ConvertToTypes()
	if err != nil {
		return &status, err
	}
	// if tx is for confChange, we should check if the client is system admin
	// get tx type according to recipient
	txType, err := types.GetTxType(common.BytesToAddress(tx.Data.Recipient).String())
	if err == nil && (txType == types.VALIDATOR || txType == types.NODE) {
		pk, err := crypto.NewPublicKey(req.PK)
		if err != nil {
			return &status, err
		}
		// create member to check if the client is system admin
		member, err := types.NewMember(pk, "")
		if err != nil {
			return &status, err
		}
		if !s.cc.CM.IsSystemAdmin(member) { // not system admin, return error
			return &status, fmt.Errorf("The client is not system admin and can't config the cluster.")
		}
	}
	err = s.cc.AddTx(tx)
	return &status, err
}
