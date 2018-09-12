package server

import (
	pb "madledger/protos"

	"golang.org/x/net/context"
)

// FetchBlock is the implementation of protos
func (s *Server) FetchBlock(ctx context.Context, req *pb.FetchBlockRequest) (*pb.Block, error) {
	block, err := s.ChannelManager.FetchBlock(req.ChannelID, req.Number)
	if err != nil {
		return nil, err
	}
	return pb.NewBlock(block)
}

// ListChannels is the implementation of protos
func (s *Server) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	return s.ChannelManager.ListChannels(req), nil
}

// AddChannel is the implementation of protos
func (s *Server) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ChannelInfo, error) {
	return s.ChannelManager.AddChannel(req)
}

// AddTx is the implementation of protos
func (s *Server) AddTx(ctx context.Context, req *pb.AddTxRequest) (*pb.TxStatus, error) {
	var status pb.TxStatus
	tx, err := req.Tx.ConvertToTypes()
	if err != nil {
		return &status, err
	}
	err = s.ChannelManager.AddTx(tx)
	return &status, err
}
