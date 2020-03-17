// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package server

import (
	"context"
	pb "madledger/protos"
)

// GetTxStatus is the implementation of protos
func (s *Server) GetTxStatus(ctx context.Context, req *pb.GetTxStatusRequest) (*pb.TxStatus, error) {
	log.Infof("in peer get tx status: channel: %s, Tx: %s", req.ChannelID, req.TxID)
	status, err := s.cm.GetTxStatus(req.ChannelID, req.TxID, true)
	log.Infof("get tx %s status of channel %s", req.TxID, req.ChannelID)
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
	log.Infof("status is %v", result)
	return result, nil
}

// ListTxHistory is the implementation of protos
// TODO: make sure the address is right and with signature
func (s *Server) ListTxHistory(ctx context.Context, req *pb.ListTxHistoryRequest) (*pb.TxHistory, error) {
	history := s.cm.GetTxHistory(req.Address)
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
