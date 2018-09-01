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
	return ConvertBlockFromTypesToPb(block)
}

// ListChannels is the implementation of protos
func (s *Server) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	return s.ChannelManager.ListChannels(req), nil
}
