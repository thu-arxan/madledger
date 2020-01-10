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
)

// GetTxType return tx type
func GetTxType(recipient string) (TxType, error) {
	if strings.Compare(recipient, CreateChannelContractAddress.String()) == 0 {
		return CREATECHANNEL, nil
	} else if strings.Compare(recipient, CfgTendermintAddress.String()) == 0 {
		return VALIDATOR, nil
	} else if strings.Compare(recipient, CfgRaftAddress.String()) == 0 {
		return NODE, nil
	} else {
		return 0, errors.New("unknown tx type")
	}
}
