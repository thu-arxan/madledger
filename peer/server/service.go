package server

import (
	"context"
	pb "madledger/protos"
)

// GetTxStatus is the implementation of protos
func (s *Server) GetTxStatus(ctx context.Context, req *pb.GetTxStatusRequest) (*pb.TxStatus, error) {
	status, err := s.ChannelManager.GetTxStatus(req.ChannelID, req.TxID)
	if err != nil {
		return &pb.TxStatus{}, nil
	}
	return &pb.TxStatus{
		Err:             status.Err,
		BlockNumber:     status.BlockNumber,
		BlockIndex:      int32(status.BlockIndex),
		Output:          status.Output,
		ContractAddress: status.ContractAddress,
	}, nil
}
