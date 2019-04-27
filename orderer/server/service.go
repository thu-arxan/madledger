package server

import (
	"errors"
	"madledger/core/types"
	pb "madledger/protos"

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
	err = s.cc.AddTx(tx)
	return &status, err
}
