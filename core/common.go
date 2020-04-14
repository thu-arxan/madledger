// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"errors"
	"madledger/common"
	"strings"
)

const (
	// GLOBALCHANNELID is the id of global channel
	GLOBALCHANNELID = "_global"
	// CONFIGCHANNELID is the id of config channel
	CONFIGCHANNELID = "_config"
	// ASSETCHANNELID is the id of account channel
	ASSETCHANNELID = "_asset"
)

var (
	// GenesisBlockPrevHash is the prev hash of genesis block
	GenesisBlockPrevHash = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

// Defines some native contracts.
var (
	// Create a channel
	CreateChannelContractAddress = common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff")
	// Config consensus cluster
	CfgConsensusAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")
	// issue
	IssueContractAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffd")
	// transfer
	TransferContractrAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffc")
	// exchange token
	TokenExchangeAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffb")
)

// GetTxType return tx type
func GetTxType(recipient string) (TxType, error) {
	if strings.Compare(recipient, CreateChannelContractAddress.String()) == 0 {
		return CREATECHANNEL, nil
	} else if strings.Compare(recipient, CfgConsensusAddress.String()) == 0 {
		return CONSENSUS, nil
	} else if strings.Compare(recipient, IssueContractAddress.String()) == 0 {
		return ISSUE, nil
	} else if strings.Compare(recipient, TransferContractrAddress.String()) == 0 {
		return TRANSFER, nil
	} else if strings.Compare(recipient, TokenExchangeAddress.String()) == 0 {
		return TOKEN, nil
	} else {
		return 0, errors.New("unknown tx type")
	}
}
