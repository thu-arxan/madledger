package channel

import (
	"encoding/binary"
	"errors"
	"madledger/blockchain"
	"madledger/common"
	"madledger/common/util"
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
	log.Infof("channel %s is starting...", m.id)
	for {
		block, err := m.fetchBlock()
		// fmt.Println("Succeed to fetch block", m.id, ":", block.Header.Number)
		if err == nil {
			m.AddBlock(block)
		} else if err.Error() == "Stop" {
			m.stopCh <- true
			return
		} else {
			log.Warnf("failed to fetch block: %d, err: %v", m.cm.GetExpect(), err)
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
	switch block.Header.ChannelID {
	case core.GLOBALCHANNELID:
		m.AddGlobalBlock(block)
		log.Infof("Add global block %d", block.Header.Number)
	case core.CONFIGCHANNELID:
		m.AddConfigBlock(block)
		log.Infof("Add config block %d", block.Header.Number)
	case core.ASSETCHANNELID:
		m.AddAssetBlock(block)
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

		wb.PutBlock(block)
		return wb.Sync()
	}

	return nil
}

// RunBlock will carry out all txs in the block.
// It will return after the block is runned.
// In the future, this will contains chains which rely on something or nothing
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
	// TODO: Gas 是不是应该判定一下如果是系统通道就不用这几个量了？

	var maxGas uint64
	maxGasByte, err := m.db.Get(util.BytesCombine([]byte(m.id), []byte("maxgas")))
	if err != nil {
		return nil, err
	}
	maxGas = uint64(binary.BigEndian.Uint64(maxGasByte))
	// var ratio uint64
	// ratioByte, err := m.db.Get(util.BytesCombine([]byte(m.id), []byte("ratio")))
	// if err != nil {
	// 	return nil, err
	// }
	// ratio = uint64(binary.BigEndian.Uint64(ratioByte))
	var gasPrice uint64
	gasPriceByte, err := m.db.Get(util.BytesCombine([]byte(m.id), []byte("gasprice")))
	if err != nil {
		return nil, err
	}
	gasPrice = uint64(binary.BigEndian.Uint64(gasPriceByte))

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
		/* TODO: gas
		用户的参数：user gas limit
		通道的参数：channel gas limit, gas price, asset token ratio
		gas limit = min (user, channel)
		获取sender的token，如果比gas limit * gas price 小，那么不能执行，直接下一个tx
		记录进入evm前的gas limit为before
		用出来之后用before-gas limit可得到具体消耗了多少gas
		然后将token -= gas * gas price，存到wb中
		*/

		log.Debug("try to find sender's token")
		gasLimit := min(tx.Data.Gas, maxGas)
		key := util.BytesCombine([]byte("token"), []byte(m.id), senderAddress.Bytes())
		tokenByte, err := m.db.Get(key)
		if err != nil {
			log.Debugf("the err is %v", err.Error())
			status.Err = err.Error()
			wb.SetTxStatus(tx, status)
			continue
		}

		tokenLeft := binary.BigEndian.Uint64(tokenByte)
		if tokenLeft < gasLimit {
			log.Info("Not enough token")
			status.Err = "Not enough token"
			wb.SetTxStatus(tx, status)
			continue
		}

		before := new(uint64)
		*before = gasLimit
		log.Infof("the gas limit is %v", &gasLimit)
		log.Infof("the before is %v", before)
		evm := evm.NewEVM(context, senderAddress, tx.Data.Payload, tx.Data.Value, gasLimit, m.db, wb)

		gasUsed := *before - gasLimit
		log.Infof("the gas cost is %v", gasUsed)
		tokenLeft -= gasUsed * gasPrice
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(tokenLeft))
		wb.Put(util.BytesCombine([]byte("token"), []byte(m.id), senderAddress.Bytes()), buf)

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

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// todo: here we should support evil orderer
func (m *Manager) fetchBlock() (*core.Block, error) {
	var lock sync.RWMutex
	var errs = make(chan error, len(m.clients))
	var blocks = make(chan *core.Block, len(m.clients))
	closed := false
	id := m.id
	expect := m.cm.GetExpect()
	for i := range m.clients {
		go func(i int) {
			block, err := m.clients[i].FetchBlock(id, expect, true)
			lock.RLock()
			defer lock.RUnlock()
			if closed {
				return
			}
			if err != nil {
				errs <- err
			} else {
				blocks <- block
			}
		}(i)
	}

	fails := 0

	for {
		defer func() {
			lock.Lock()
			if !closed {
				close(errs)
				close(blocks)
				closed = true
			}
			lock.Unlock()
		}()
		// log.Infof("get %s %d, closed: %t, fail: %d, %d", id, except, closed, fails, len(m.clients))
		select {
		case block := <-blocks:
			return block, nil
		case err := <-errs:
			fails++
			if fails == len(m.clients) {
				return nil, err
			}
		case <-m.signalCh:
			return nil, errors.New("Stop")
		}
	}
}
