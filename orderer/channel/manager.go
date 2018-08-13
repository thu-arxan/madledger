package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/core"
	"madledger/orderer/db"

	"github.com/rs/zerolog/log"
)

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
	// db is the database
	db *db.DB
	// chain manager
	cm *blockchain.Manager
}

// NewManager is the constructor of Manager
// TODO: many things is not done yet
func NewManager(id, dir string, db *db.DB) (*Manager, error) {
	log.Info().Msg(dir)
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID: id,
		db: db,
		cm: cm,
	}, nil
}

// Start starts the channel
// TODO
func (manager *Manager) Start() {
	log.Info().Msgf("Channel %s is starting")
}

// FetchBlock return the block if exist
// TODO
func (manager *Manager) FetchBlock(num uint64) (*core.Block, error) {
	return nil, errors.New("Not implementation yet")
}
