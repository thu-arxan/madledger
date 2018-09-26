package db

import (
	cc "madledger/blockchain/config"
	"madledger/core/types"
)

// DB is the interface of db
type DB interface {
	// ListChannel list all channels
	ListChannel() []string
	HasChannel(id string) bool
	UpdateChannel(id string, profile *cc.Profile) error
	// AddBlock will records all txs in the block to get rid of duplicated txs
	AddBlock(block *types.Block) error
	HasTx(tx *types.Tx) bool
	IsMember(channelID string, member *types.Member) bool
	IsAdmin(channelID string, member *types.Member) bool
	Close() error
}
