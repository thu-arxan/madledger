package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/common"
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
	identity *types.Member
	// id is the id of channel
	id string
	// db is the database
	db db.DB
	// chain manager
	cm          *blockchain.Manager
	client      *orderer.Client
	coordinator *Coordinator
}

// NewManager is the constructor of Manager
func NewManager(id, dir string, identity *types.Member, db db.DB, client *orderer.Client, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		identity:    identity,
		id:          id,
		db:          db,
		cm:          cm,
		client:      client,
		coordinator: coordinator,
	}, nil
}

// Start start the manager.
func (m *Manager) Start() {
	log.Infof("%s is starting...", m.id)
	for {
		block, err := m.fetchBlock()
		// fmt.Println("Succeed to fetch block", m.id, ":", block.Header.Number)
		if err == nil {
			m.AddBlock(block)
		}
	}
}

// Stop will stop the manager
// TODO: find a good way to stop
func (m *Manager) Stop() {
}

// AddBlock add a block
func (m *Manager) AddBlock(block *types.Block) error {
	// add into the blockchain
	err := m.cm.AddBlock(block)
	if err != nil {
		return err
	}
	switch block.Header.ChannelID {
	case types.GLOBALCHANNELID:
		m.AddGlobalBlock(block)
		log.Infof("Add global block %d", block.Header.Number)
	case types.CONFIGCHANNELID:
		m.AddConfigBlock(block)
	default:
		for {
			if m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
				log.Infof("Run block %s:%d", m.id, block.Header.Number)
				m.RunBlock(block.Header.Number)
				return nil
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

// RunBlock will carry out all txs in the block.
// It will return after the block is runned.
// In the future, this will contains chains which rely on something or nothing
// TODO: transfer is not implementation yet
func (m *Manager) RunBlock(num uint64) error {
	block, err := m.cm.GetBlock(num)
	if err != nil {
		return err
	}
	context := evm.NewContext(block)
	for i, tx := range block.Transactions {
		senderAddress, err := tx.GetSender()
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: num,
			BlockIndex:  i,
			Output:      nil,
		}
		if err != nil {
			status.Err = err.Error()
			m.db.SetTxStatus(tx, status)
			continue
		}
		receiverAddress := tx.GetReceiver()
		log.Infof("The address of receiver is %s", receiverAddress.String())
		sender, err := m.db.GetAccount(senderAddress)
		if err != nil {
			status.Err = err.Error()
			m.db.SetTxStatus(tx, status)
			continue
		}

		evm := evm.NewEVM(*context, senderAddress, m.db)
		if receiverAddress.String() != common.ZeroAddress.String() {
			// log.Info("This is a normal call")
			// if the length of payload is not zero, this is a contract call
			if len(tx.Data.Payload) != 0 && !m.db.AccountExist(receiverAddress) {
				status.Err = "Invalid Address"
				m.db.SetTxStatus(tx, status)
				continue
			}

			receiver, err := m.db.GetAccount(receiverAddress)
			if err != nil {
				status.Err = err.Error()
				m.db.SetTxStatus(tx, status)
				continue
			}
			output, err := evm.Call(sender, receiver, receiver.GetCode(), tx.Data.Payload, 0)
			status.Output = output
			if err != nil {
				status.Err = err.Error()
			}
			m.db.SetTxStatus(tx, status)
		} else {
			log.Info("This is a create call")
			output, addr, err := evm.Create(sender, tx.Data.Payload, []byte{}, 0)
			status.Output = output
			status.ContractAddress = addr.String()
			if err != nil {
				status.Err = err.Error()
			}
			m.db.SetTxStatus(tx, status)
		}
	}
	return errors.New("Not implementation yet")
}

func (m *Manager) fetchBlock() (*types.Block, error) {
	return m.client.FetchBlock(m.id, m.cm.GetExcept(), true)
}
