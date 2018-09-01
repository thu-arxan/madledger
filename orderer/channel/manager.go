package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/core/types"
	"madledger/orderer/db"

	"github.com/rs/zerolog/log"
)

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
	// db is the database
	db db.DB
	// chain manager
	cm *blockchain.Manager
}

// NewManager is the constructor of Manager
// TODO: many things is not done yet
func NewManager(id, dir string, db db.DB) (*Manager, error) {
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

// HasGenesisBlock return if the channel has a genesis block
func (manager *Manager) HasGenesisBlock() bool {
	return manager.cm.HasGenesisBlock()
}

// GetBlock return the block of num
func (manager *Manager) GetBlock(num uint64) (*types.Block, error) {
	return manager.cm.GetBlock(num)
}

// AddBlock add a block
// TODO: check conflict and update db
func (manager *Manager) AddBlock(block *types.Block) error {
	var err error
	err = manager.cm.AddBlock(block)
	if err != nil {
		return err
	}
	switch manager.ID {
	case types.CONFIGCHANNELID:
		return manager.AddConfigBlock(block)
	case types.GLOBALCHANNELID:
		return nil
	default:
		// todo
		return nil
	}
}

// Start starts the channel
// TODO
func (manager *Manager) Start() {
	log.Info().Msgf("Channel %s is starting")
}

// FetchBlock return the block if exist
// TODO
func (manager *Manager) FetchBlock(num uint64) (*types.Block, error) {
	return nil, errors.New("Not implementation yet")
}
