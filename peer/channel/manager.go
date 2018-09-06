package channel

import (
	"madledger/blockchain"
	"madledger/core/types"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	// id is the id of channel
	id string
	// db is the database
	db db.DB
	// chain manager
	cm            *blockchain.Manager
	init          bool
	stop          chan bool
	ordererClient *orderer.Client
}

// NewManager is the constructor of Manager
func NewManager(id, dir string, db db.DB, ordererClient *orderer.Client) (*Manager, error) {
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		id:            id,
		db:            db,
		cm:            cm,
		init:          false,
		stop:          make(chan bool),
		ordererClient: ordererClient,
	}, nil
}

// Start start the manager.
// The manager will try to fetch a block every 500ms,
// but remember this is a implementation which is very bad.
// It should be replaced as soon as possible.
// TODO:
func (m *Manager) Start() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	log.Infof("%s is starting...", m.id)
	for {
		select {
		case <-ticker.C:
			// log.Infof("%s ticker", m.id)
			block, err := m.fetchBlock()
			if err == nil {
				m.cm.AddBlock(block)
			}
		}
	}
}

func (m *Manager) fetchBlock() (*types.Block, error) {
	return m.ordererClient.FetchBlock(m.id, m.cm.GetExcept())
}
