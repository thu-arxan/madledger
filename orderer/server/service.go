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
	"errors"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
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
	tx, err := req.GetTx().ToCore()
	if err != nil {
		return nil, err
	}
	if !tx.Verify() {
		return nil, errors.New("The tx is not a valid tx")
	}
	if tx.GetReceiver().String() != core.CreateChannelContractAddress.String() {
		return nil, errors.New("The receiver of the tx is not the valid contract address")
	}
	_, err = s.cc.CreateChannel(tx)
	if err != nil {
		return nil, err
	}
	return &pb.ChannelInfo{}, nil
}

// AddTx is the implementation of protos
// TODO: Tx must be verified
func (s *Server) AddTx(ctx context.Context, req *pb.AddTxRequest) (*pb.TxStatus, error) {
	var status pb.TxStatus
	if req.Tx == nil {
		return &status, errors.New("tx should not be nil")
	}
	tx, err := req.Tx.ToCore()
	if err != nil {
		return &status, err
	}
	if !tx.Verify() {
		return &status, errors.New("tx is not formatted")
	}
	// if tx is for confChange, we should check if the client is system admin
	// get tx type according to recipient
	txType, err := core.GetTxType(common.BytesToAddress(tx.Data.Recipient).String())
	if err == nil && txType == core.CONSENSUS {
		pk, err := crypto.NewPublicKey(req.Tx.Data.Sig.PK, req.Tx.Data.Sig.Algo)
		if err != nil {
			return &status, err
		}
		// create member to check if the client is system admin
		member, err := core.NewMember(pk, "")
		if err != nil {
			return &status, err
		}
		if !s.cc.CM.IsSystemAdmin(member) { // not system admin, return error
			return &status, errors.New("The client is not system admin and can't config the cluster")
		}
	}

	err = s.cc.AddTx(tx)
	return &status, err
}

// GetAccountInfo is the implementation of protos
func (s *Server) GetAccountInfo(ctx context.Context, req *pb.GetAccountInfoRequest) (*pb.AccountInfo, error) {
	var info pb.AccountInfo
	address := common.BytesToAddress(req.Address)
	account, err := s.cc.AM.GetAccount(address)
	if err != nil {
		return &info, err
	}
	info.Balance = account.GetBalance()
	return &info, err
}
