package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/common"
	"madledger/core/types"
	"madledger/executor/evm"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	signalCh chan bool
	stopCh   chan bool
	identity *types.Member
	// id is the id of channel
	id string
	// db is the database
	db db.DB
	// chain manager
	cm          *blockchain.Manager
	clients     []*orderer.Client
	coordinator *Coordinator
}

// NewManager is the constructor of Manager
func NewManager(id, dir string, identity *types.Member, db db.DB, clients []*orderer.Client, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		signalCh:    make(chan bool, 1),
		stopCh:      make(chan bool, 1),
		identity:    identity,
		id:          id,
		db:          db,
		cm:          cm,
		clients:     clients,
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
		} else if err.Error() == "Stop" {
			m.stopCh <- true
			return
		}
	}
}

// Stop will stop the manager
// TODO: find a good way to stop
func (m *Manager) Stop() {
	m.signalCh <- true
	<-m.stopCh
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
		log.Infof("Add config block %d", block.Header.Number)
	default:
		for {
			if m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
				log.Infof("Run block %s: %d", m.id, block.Header.Number)
				wb, err := m.RunBlock(block.Header.Number)
				if err != nil {
					return err
				}
				batch := wb.GetBatch()
				err = m.db.SyncWriteBatch(batch)
				if err != nil {
					return err
				}
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
func (m *Manager) RunBlock(num uint64) (db.WriteBatch, error) {
	block, err := m.cm.GetBlock(num)
	if err != nil {
		return nil, err
	}
	context := evm.NewContext(block)
	wb := m.db.NewWriteBatch()
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
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
			continue
		}
		receiverAddress := tx.GetReceiver()
		log.Infof("The address of receiver is %s", receiverAddress.String())
		sender, err := m.db.GetAccount(senderAddress)
		if err != nil {
			status.Err = err.Error()
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == types.CfgTendermintAddress.String() {
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == types.CfgRaftAddress.String() {
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
			continue
		}

		evm := evm.NewEVM(*context, senderAddress, m.db)
		if receiverAddress.String() != common.ZeroAddress.String() {
			// log.Info("This is a normal call")
			// if the length of payload is not zero, this is a contract call
			if len(tx.Data.Payload) != 0 && !m.db.AccountExist(receiverAddress) {
				status.Err = "Invalid Address"
				//m.db.SetTxStatus(tx, status)
				wb.SetTxStatus(tx, status)
				continue
			}

			receiver, err := m.db.GetAccount(receiverAddress)
			if err != nil {
				status.Err = err.Error()
				//m.db.SetTxStatus(tx, status)
				wb.SetTxStatus(tx, status)
				continue
			}
			output, err := evm.Call(sender, receiver, receiver.GetCode(), tx.Data.Payload, 0, wb)
			status.Output = output
			if err != nil {
				status.Err = err.Error()
			}
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
		} else {
			log.Info("This is a create call")
			output, addr, err := evm.Create(sender, tx.Data.Payload, []byte{}, 0, wb)
			status.Output = output
			status.ContractAddress = addr.String()
			if err != nil {
				status.Err = err.Error()
			}
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
		}
	}
	return wb, nil
}

// todo: here we should support evil orderer
func (m *Manager) fetchBlock() (*types.Block, error) {
	var lock sync.Mutex
	var ch = make(chan bool, 1)
	var errs = make([]error, len(m.clients))
	var blocks = make([]*types.Block, len(m.clients))
	id := m.id
	except := m.cm.GetExcept()
	for i := range m.clients {
		go func(i int) {
			block, err := m.clients[i].FetchBlock(id, except, true)
			if err != nil {
				errs[i] = err
				lock.Lock()
				defer lock.Unlock()
				for i := range errs {
					if errs[i] == nil {
						return
					}
				}
				ch <- false
			} else {
				blocks[i] = block
				ch <- true
			}
		}(i)
	}

	select {
	case ok := <-ch:
		if ok {
			for i := range blocks {
				if blocks[i] != nil {
					return blocks[i], nil
				}
			}
		}
		return nil, errs[0]
	case <-m.signalCh:
		return nil, errors.New("Stop")
	}
}
