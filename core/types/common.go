package types

import "madledger/common"

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
	CfgTendermintAddress=common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")
	// Config the raft cluster
	CfgRaftAddress=common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffd")
)
