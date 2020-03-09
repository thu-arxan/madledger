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
	// Config the tendermint cluster
	CfgTendermintAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")
	// Config the raft cluster
	CfgRaftAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffd")
	// issue
	IssueContractAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffc")
	// transfer
	TransferContractrAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffb")
)

// GetTxType return tx type
func GetTxType(recipient string) (TxType, error) {
	if strings.Compare(recipient, CreateChannelContractAddress.String()) == 0 {
		return CREATECHANNEL, nil
	} else if strings.Compare(recipient, CfgTendermintAddress.String()) == 0 {
		return VALIDATOR, nil
	} else if strings.Compare(recipient, CfgRaftAddress.String()) == 0 {
		return NODE, nil
	} else if strings.Compare(recipient, IssueContractAddress.String()) == 0 {
		return ISSUE, nil
	} else if strings.Compare(recipient, TransferContractrAddress.String()) == 0 {
		return TRANSFER, nil
	} else {
		return 0, errors.New("unknown tx type")
	}
}
