package channel

import (
	"errors"
	"fmt"
	"madledger/blockchain"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core/types"
	"madledger/orderer/db"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "orderer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
	// db is the database
	db db.DB
	// chain manager
	cm *blockchain.Manager
	// consensus block chan
	cbc         chan consensus.Block
	init        bool
	stop        chan bool
	hub         *event.Hub
	coordinator *Coordinator
}

// NewManager is the constructor of Manager
func NewManager(id string, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, fmt.Sprintf("%s/%s", coordinator.chainCfg.Path, id))
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID:          id,
		db:          coordinator.db,
		cm:          cm,
		cbc:         make(chan consensus.Block),
		init:        false,
		stop:        make(chan bool),
		hub:         event.NewHub(),
		coordinator: coordinator,
	}, nil
}

// Start starts the channel
func (manager *Manager) Start() {
	log.Infof("Channel %s is starting", manager.ID)
	// manager.coordinator.Consensus.SyncBlocks(manager.ID, &(manager.consensusBlockChan))
	manager.init = true
	go manager.syncBlock()
	for {
		select {
		case cb := <-manager.cbc:
			log.Infof("Receive block %d from consunsus\n", cb.GetNumber())
			// todo: if a tx is duplicated and it was added into consensus block succeed, then it will never receive response
			txs := removeDuplicateTxs(manager.db, GetTxsFromConsensusBlock(cb))
			if len(txs) != 0 {
				prevBlock := manager.cm.GetPrevBlock()
				var block *types.Block
				if prevBlock == nil {
					block = types.NewBlock(manager.ID, 0, types.GenesisBlockPrevHash, txs)
					log.Infof("Channel %s create new block %d, hash is %s", manager.ID, 0, util.Hex(block.Hash().Bytes()))
				} else {
					block = types.NewBlock(manager.ID, prevBlock.Header.Number+1, prevBlock.Hash().Bytes(), txs)
					log.Infof("Channel %s create new block %d, hash is %s", manager.ID, prevBlock.Header.Number+1, util.Hex(block.Hash().Bytes()))
				}
				// If the channel is the global channel, the block is finished.
				// else send a tx to the global channel
				if manager.ID != types.GLOBALCHANNELID {
					tx := types.NewGlobalTx(manager.ID, block.Header.Number, block.Hash())
					err := manager.coordinator.GM.AddTx(tx)
					if err != nil {
						log.Fatalf("Channel %s failed to add tx into global channel because %s", manager.ID, err)
						return
					}
				}
				err := manager.AddBlock(block)
				if err != nil {
					log.Fatalf("Channel %s failed to run because of %s", manager.ID, err)
					return
				}
				log.Infof("Channel %s has %d block now", manager.ID, block.Header.Number+1)
				manager.hub.Done(string(block.Header.Number), nil)
				for _, tx := range block.Transactions {
					manager.hub.Done(util.Hex(tx.Hash()), nil)
				}
			}
		case <-manager.stop:
			manager.init = false
			return
		}
	}
}

func (manager *Manager) syncBlock() {
	var num uint64 = 1
	for {
		cb, err := manager.coordinator.Consensus.GetBlock(manager.ID, num, true)
		if err != nil {
			fmt.Println(err)
			continue
		}
		num++
		go func() {
			manager.cbc <- cb
		}()
	}
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
func (manager *Manager) AddBlock(block *types.Block) error {
	// first update db
	if err := manager.db.AddBlock(block); err != nil {
		return err
	}
	if err := manager.cm.AddBlock(block); err != nil {
		return err
	}
	// check is there is any need to update local state of orderer
	switch manager.ID {
	case types.CONFIGCHANNELID:
		return manager.AddConfigBlock(block)
	case types.GLOBALCHANNELID:
		return nil
	default:
		return nil
	}
}

// GetBlockSize return the size of blocks
func (manager *Manager) GetBlockSize() uint64 {
	return manager.cm.GetExcept()
}

// AddTx try to add a tx
func (manager *Manager) AddTx(tx *types.Tx) error {
	if manager.db.HasTx(tx) {
		return errors.New("The tx exist in the blockchain aleardy")
	}

	txBytes, err := tx.Bytes()
	if err != nil {
		return err
	}

	err = manager.coordinator.Consensus.AddTx(manager.ID, txBytes)
	if err != nil {
		return err
	}
	// Note: The reason why we must do this is because we must make sure we return the result after we store the block
	// However, we may find a better way to do this if we allow there are more interactive between the consensus and orderer.
	result := manager.hub.Watch(util.Hex(tx.Hash()), nil)

	return result.Err
}

// Stop stop the manager
func (manager *Manager) Stop() {
	if manager.init {
		manager.stop <- true
		for manager.init {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// FetchBlock return the block if exist
func (manager *Manager) FetchBlock(num uint64) (*types.Block, error) {
	return manager.cm.GetBlock(num)
}

// IsMember return if the member belongs to the channel
func (manager *Manager) IsMember(member *types.Member) bool {
	return manager.db.IsMember(manager.ID, member)
}

// IsAdmin return if the member is the admin of the channel
func (manager *Manager) IsAdmin(member *types.Member) bool {
	return manager.db.IsAdmin(manager.ID, member)
}

// FetchBlockAsync will fetch book async.
// TODO: fix the thread unsafety
func (manager *Manager) FetchBlockAsync(num uint64) (*types.Block, error) {
	if manager.cm.GetExcept() <= num {
		manager.hub.Watch(string(num), nil)
	}

	block, err := manager.cm.GetBlock(num)
	if err == nil {
		return block, err
	}
	return nil, err
}

// removeDuplicateTxs will remove tx which exists in the previous blocks
func removeDuplicateTxs(db db.DB, txs []*types.Tx) []*types.Tx {
	var unduplicateTxs []*types.Tx
	for _, tx := range txs {
		if !util.Contain(unduplicateTxs, tx) {
			if !db.HasTx(tx) {
				unduplicateTxs = append(unduplicateTxs, tx)
			}
		}
	}
	return unduplicateTxs
}
