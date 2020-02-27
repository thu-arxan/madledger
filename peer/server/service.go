package server

import (
	"context"
	pb "madledger/protos"
)

// GetTxStatus is the implementation of protos
func (s *Server) GetTxStatus(ctx context.Context, req *pb.GetTxStatusRequest) (*pb.TxStatus, error) {
	status, err := s.ChannelManager.GetTxStatus(req.ChannelID, req.TxID, true)
	log.Debugf("get tx %s status of channel %s", req.TxID, req.ChannelID)
	if err != nil {
		return &pb.TxStatus{}, err
	}
	result := &pb.TxStatus{
		Err:             status.Err,
		BlockNumber:     status.BlockNumber,
		BlockIndex:      int32(status.BlockIndex),
		Output:          status.Output,
		ContractAddress: status.ContractAddress,
	}
	return result, nil
}

// ListTxHistory is the implementation of protos
// TODO: make sure the address is right and with signature
func (s *Server) ListTxHistory(ctx context.Context, req *pb.ListTxHistoryRequest) (*pb.TxHistory, error) {
	history := s.ChannelManager.ListTxHistory(req.Address)
	var pbHistory = make(map[string]*pb.StringList)
	for channelID, ids := range history {
		value := new(pb.StringList)
		for _, id := range ids {
			value.Value = append(value.Value, id)
		}
		pbHistory[channelID] = value
	}
	return &pb.TxHistory{
		Txs: pbHistory,
	}, nil
}
