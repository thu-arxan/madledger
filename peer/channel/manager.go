package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/executor/evm"
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
			block, err := m.fetchBlock()
			if err == nil {
				// m.cm.AddBlock(block)
				m.AddBlock(block)
			}
		}
	}
}

// AddBlock add a block
func (m *Manager) AddBlock(block *types.Block) error {
	// add into the blockchain
	m.cm.AddBlock(block)
	switch block.Header.ChannelID {
	case types.GLOBALCHANNELID:
		m.AddGlobalBlock(block)
	case types.CONFIGCHANNELID:
		// todo
	default:
		// do nothing now
		m.RunBlock(block.Header.Number)
	}
	return nil
}

// RunBlock will carry out all txs in the block.
// It will return after the block is runned.
// In the future, this will contains chains which rely on something or nothing
func (m *Manager) RunBlock(num uint64) error {
	block, err := m.cm.GetBlock(num)
	if err != nil {
		return err
	}
	context := evm.NewContext(block)
	for _, tx := range block.Transactions {
		senderAddress, err := tx.GetSender()
		if err == nil {
			receiverAddress := tx.GetReceiver()
			log.Infof("The address of receiver is %s", receiverAddress.String())
			sender, err := m.db.GetAccount(senderAddress)
			if err != nil {
				continue
			}
			receiver, err := m.db.GetAccount(receiverAddress)
			if err != nil {
				continue
			}
			evm := evm.NewEVM(*context, senderAddress, m.db)
			log.Infof("The address of receiver is %s", receiver.GetAddress().String())
			if receiver.GetAddress().String() != common.ZeroAddress.String() {
				log.Info("This is a normal call")
				log.Info(util.Hex(receiver.GetCode()))
				output, err := evm.Call(sender, receiver, receiver.GetCode(), tx.Data.Payload, 0)
				if err != nil {
					log.Error(err)
				} else {
					log.Info(output)
				}
			} else {
				log.Info("This is a create call")
				_, addr, err := evm.Create(sender, tx.Data.Payload, []byte{}, 0)
				if err != nil {
					log.Error(err)
				}
				log.Infof("Get address %s", addr.String())
			}
		}

	}
	return errors.New("Not implementation yet")
}

func (m *Manager) fetchBlock() (*types.Block, error) {
	return m.ordererClient.FetchBlock(m.id, m.cm.GetExcept())
}
