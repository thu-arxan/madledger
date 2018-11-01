package channel

import (
	"errors"
	"fmt"
	"madledger/blockchain"
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
	cm                 *blockchain.Manager
	consensusBlockChan chan consensus.Block
	init               bool
	stop               chan bool
	notify             *notifyPool
	coordinator        *Coordinator
}

// NewManager is the constructor of Manager
func NewManager(id string, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, fmt.Sprintf("%s/%s", coordinator.chainCfg.Path, id))
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID:                 id,
		db:                 coordinator.db,
		cm:                 cm,
		consensusBlockChan: make(chan consensus.Block),
		init:               false,
		stop:               make(chan bool),
		notify:             newNotifyPool(),
		coordinator:        coordinator,
	}, nil
}

// Start starts the channel
func (manager *Manager) Start() {
	log.Infof("Channel %s is starting", manager.ID)
	// manager.consensus = consensus
	// manager.globalManager = globalManager
	manager.coordinator.Consensus.SyncBlocks(manager.ID, &(manager.consensusBlockChan))
	manager.init = true
	for {
		select {
		case cb := <-manager.consensusBlockChan:
			// However, this does not means that the block is added succeed, because it may need to be added into global channel.
			// So here are many things need to be done.
			txs := GetTxsFromConsensusBlock(cb)
			var unduplicateTxs []*types.Tx
			for _, tx := range txs {
				if !util.Contain(unduplicateTxs, tx) {
					if !manager.db.HasTx(tx) {
						unduplicateTxs = append(unduplicateTxs, tx)
					}
				}
			}
			if len(unduplicateTxs) == 0 {
				return
			}
			prevBlock := manager.cm.GetPrevBlock()
			var block *types.Block
			if prevBlock == nil {
				block = types.NewBlock(manager.ID, 0, types.GenesisBlockPrevHash, txs)
				log.Infof("Channel %s create new block %d, hash is %s", manager.ID, 0, util.Hex(block.Hash().Bytes()))
			} else {
				block = types.NewBlock(manager.ID, prevBlock.Header.Number+1, prevBlock.Hash().Bytes(), txs)
				log.Infof("Channel %s create new block %d, hash is %s", manager.ID, prevBlock.Header.Number+1, util.Hex(block.Hash().Bytes()))
			}
			// then if the channel is the global channel, the block is finished.
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
			manager.notify.addBlock(block)

		case <-manager.stop:
			manager.init = false
			return
		}
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
// TODO: check conflict and update db
func (manager *Manager) AddBlock(block *types.Block) error {
	var err error
	// first update db
	manager.db.AddBlock(block)
	// if err != nil {
	// 	return err
	// }
	err = manager.cm.AddBlock(block)
	if err != nil {
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
	// First check if the tx exists aleaydy, if true return error right away
	if manager.db.HasTx(tx) {
		return errors.New("The tx exist in the blockchain aleardy")
	}
	txBytes, err := tx.Bytes()
	if err != nil {
		return err
	}
	// first register a notify event
	// Note: The reason why we must do this is because we must make sure we return the result after we store the block
	// However, we may find a better way to do this if we allow there are more interactive between the consensus and orderer.
	var errChan = make(chan error)
	var hash = util.Hex(tx.Hash())
	err = manager.notify.addTxNotify(hash, &errChan)
	if err != nil {
		return err
	}

	defer manager.notify.deleteTxNotify(hash)
	err = manager.coordinator.Consensus.AddTx(manager.ID, txBytes)
	if err != nil {
		return err
	}
	err = <-errChan
	return err
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
// However, it would be better if using a pool rather than using while.
// TODO: fix the thread unsafety
func (manager *Manager) FetchBlockAsync(num uint64) (*types.Block, error) {
	var b = make(chan bool)
	if manager.cm.GetExcept() > num {
		block, err := manager.cm.GetBlock(num)
		if err == nil {
			return block, err
		}
		return nil, err
	}
	manager.notify.addBlockNotify(num, &b)
	<-b
	block, err := manager.cm.GetBlock(num)
	if err == nil {
		return block, err
	}
	return nil, err
}
