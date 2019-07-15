package db

import (
	cc "madledger/blockchain/config"
	"madledger/core/types"
)

// DB is the interface of db, and it is the implementation of DB on orderer/.tendermint/.glue
type DB interface {
	ListChannel() []string
	HasChannel(id string) bool
	UpdateChannel(id string, profile *cc.Profile) error
	// AddBlock will records all txs in the block to get rid of duplicated txs
	AddBlock(block *types.Block) error
	HasTx(tx *types.Tx) bool
	IsMember(channelID string, member *types.Member) bool
	IsAdmin(channelID string, member *types.Member) bool
	// WatchChannel provide a way to spy channel change. Now it mainly used to
	// spy channel create operation.
	WatchChannel(channelID string)
	Close() error
	UpdateSystemAdmin(profile *cc.Profile) error
	IsSystemAdmin(member *types.Member) bool
}
