package channel

import (
	"madledger/orderer/db"

	"github.com/rs/zerolog/log"
)

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
	// db is the database
	db *db.DB
}

// NewManager is the constructor of Manager
// TODO: many things is not done yet
func NewManager(id string, db *db.DB) (*Manager, error) {
	return &Manager{
		ID: id,
		db: db,
	}, nil
}

// Start starts the channel
// TODO
func (manager *Manager) Start() {
	log.Info().Msgf("Channel %s is starting")
}
