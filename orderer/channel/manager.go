package channel

import (
	"encoding/json"
	"madledger/blockchain"
	"madledger/consensus"
	"madledger/core/types"
	"madledger/orderer/db"
	"time"
	"transaction_service/util"

	"github.com/rs/zerolog/log"
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
	consensus          consensus.Consensus
	notify             *notifyPool
}

// NewManager is the constructor of Manager
// TODO: many things is not done yet
func NewManager(id, dir string, db db.DB) (*Manager, error) {
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID:                 id,
		db:                 db,
		cm:                 cm,
		consensusBlockChan: make(chan consensus.Block),
		init:               false,
		stop:               make(chan bool),
		notify:             newNotifyPool(),
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
	// check is there is any need to update local state of orderer
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

// GetBlockSize return the size of blocks
func (manager *Manager) GetBlockSize() uint64 {
	return manager.cm.GetExcept()
}

// Start starts the channel
// TODO: many things to be done
func (manager *Manager) Start(consensus consensus.Consensus) {
	log.Info().Msgf("Channel %s is starting", manager.ID)
	manager.consensus = consensus
	consensus.SyncBlocks(manager.ID, &(manager.consensusBlockChan))
	manager.init = true
	select {
	case cb := <-manager.consensusBlockChan:
		log.Info().Msgf("Channel %s receive a block %d", manager.ID, cb.GetNumber())
		// However, this does not means that the block is added succeed, because it may need to be added into global channel.
		// So here are many things need to be done.
		txs := GetTxsFromConsensusBlock(cb)
		prevBlock := manager.cm.GetPrevBlock()
		var block *types.Block
		if prevBlock == nil {
			block = types.NewBlock(manager.ID, 0, nil, txs)
		} else {
			block = types.NewBlock(manager.ID, prevBlock.Header.Number+1, prevBlock.Hash().Bytes(), txs)
		}
		// then if the channel is the global channel, the block is finished.
		// else send a tx to the global channel
		if manager.ID != types.GLOBALCHANNELID {
			tx := types.NewGlobalTx(manager.ID, block.Header.Number, block.Hash())
			txBytes, _ := json.Marshal(tx)
			err := consensus.AddTx(types.GLOBALCHANNELID, txBytes)
			if err != nil {
				log.Fatal().Msgf("Channel %s failed to run", manager.ID)
				return
			}
		}
		err := manager.AddBlock(block)
		if err != nil {
			log.Fatal().Msgf("Channel %s failed to run", manager.ID)
			return
		}
		log.Info().Msgf("Channel %s has %d block now", manager.ID, block.Header.Number+1)
		manager.notify.addBlock(block)

	case <-manager.stop:
		manager.init = false
		return
	}
}

// AddTx try to add a tx
func (manager *Manager) AddTx(tx *types.Tx) error {
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	// first register a notify event
	var errChan = make(chan error)
	var hash = util.Hex(tx.Hash())
	err = manager.notify.addNotify(hash, &errChan)
	if err != nil {
		return err
	}

	defer manager.notify.deleteNotify(hash)
	err = manager.consensus.AddTx(manager.ID, txBytes)
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
	// return nil, errors.New("Not implementation yet")
	return manager.cm.GetBlock(num)
}
