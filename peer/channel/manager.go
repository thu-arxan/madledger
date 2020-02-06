package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/common"
	"madledger/core"

	"madledger/executor/evm"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	signalCh chan bool
	stopCh   chan bool
	identity *core.Member
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
func NewManager(id, dir string, identity *core.Member, db db.DB, clients []*orderer.Client, coordinator *Coordinator) (*Manager, error) {
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
func (m *Manager) AddBlock(block *core.Block) error {
	// add into the blockchain
	err := m.cm.AddBlock(block)
	if err != nil {
		return err
	}
	if err := m.db.PutBlock(block); err != nil {
		return err
	}
	switch block.Header.ChannelID {
	case core.GLOBALCHANNELID:
		m.AddGlobalBlock(block)
		log.Infof("Add global block %d", block.Header.Number)
	case core.CONFIGCHANNELID:
		m.AddConfigBlock(block)
		log.Infof("Add config block %d", block.Header.Number)
	case core.ACCOUNTCHANNELID:
		m.AddAccountBlock(block)
		log.Infof("Add account block %d", block.Header.Number)
	default:
		if !m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			m.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		log.Infof("Run block %s: %d", m.id, block.Header.Number)
		wb, err := m.RunBlock(block)
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

	return nil
}

// RunBlock will carry out all txs in the block.
// It will return after the block is runned.
// In the future, this will contains chains which rely on something or nothing
// TODO: transfer is not implementation yet
func (m *Manager) RunBlock(block *core.Block) (db.WriteBatch, error) {
	wb := m.db.NewWriteBatch()
	context := evm.NewContext(block, m.db, wb)
	defer context.BlockFinalize()
	// first parallel get sender to speed up
	threadSize := runtime.NumCPU()
	if threadSize < 2 {
		threadSize = 2
	}
	var ch = make(chan bool, threadSize)
	var wg sync.WaitGroup
	for i := range block.Transactions {
		wg.Add(1)
		tx := block.Transactions[i]
		ch <- true
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()
			tx.GetSender()
		}()
	}
	wg.Wait()

	for i, tx := range block.Transactions {
		senderAddress, err := tx.GetSender()
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: block.Header.Number,
			BlockIndex:  i,
			Output:      nil,
		}
		if err != nil {
			status.Err = err.Error()
			wb.SetTxStatus(tx, status)
			continue
		}
		receiverAddress := tx.GetReceiver()

		sender, err := m.db.GetAccount(senderAddress)
		if err != nil {
			status.Err = err.Error()
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == core.CfgTendermintAddress.String() {
			wb.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == core.CfgRaftAddress.String() {
			wb.SetTxStatus(tx, status)
			continue
		}
		log.Debugf(" %v, %v, %v", sender, context, tx)
		gas := uint64(10000000)
		evm := evm.NewEVM(context, senderAddress, tx.Data.Payload, tx.Data.Value, gas, m.db, wb)
		if receiverAddress.String() != common.ZeroAddress.String() {
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
			output, err := evm.Call(sender, receiver, receiver.GetCode())
			status.Output = output
			if err != nil {
				status.Err = err.Error()
			}
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
		} else {
			output, addr, err := evm.Create(sender)
			status.Output = output
			status.ContractAddress = addr.String()
			if err != nil {
				status.Err = err.Error()
			}
			//m.db.SetTxStatus(tx, status)
			wb.SetTxStatus(tx, status)
		}
	}
	// wb.PersistLog([]byte(fmt.Sprintf("block_log_%d", block.GetNumber())))
	return wb, nil
}

// todo: here we should support evil orderer
func (m *Manager) fetchBlock() (*core.Block, error) {
	var lock sync.Mutex
	var ch = make(chan bool, 1)
	var errs = make([]error, len(m.clients))
	var blocks = make([]*core.Block, len(m.clients))
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
