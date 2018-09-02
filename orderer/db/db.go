package db

import (
	cc "madledger/blockchain/config"
)

// DB is the interface of db
type DB interface {
	// ListChannel list all channels
	ListChannel() []string
	HasChannel(id string) bool
	UpdateChannel(id string, profile cc.Profile) error
	Close() error
}
