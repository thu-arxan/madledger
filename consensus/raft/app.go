package raft

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"sort"
	"sync"
	"sync/atomic"
)

// App is the application
type App struct {
	lock   sync.Mutex
	cfg    *EraftConfig
	status int32 // only Running and Stopped

	blocks  map[uint64]*HybridBlock
	blockCh chan *HybridBlock
	hub     *event.Hub
	// minBlock is the min block number that the blockchain system needed
	minBlock uint64
	db       *DB
}

// NewApp is the constructor of App
func NewApp(cfg *EraftConfig) (*App, error) {
	return &App{
		cfg:     cfg,
		blocks:  make(map[uint64]*HybridBlock),
		blockCh: make(chan *HybridBlock, 2048),
		hub:     event.NewHub(),
		status:  Stopped,
	}, nil
}

// Start start the app
func (a *App) Start() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.status != Stopped {
		return errors.New("The app is not stopped")
	}

	db, err := NewDB(a.cfg.dbDir)
	if err != nil {
		return err
	}
	a.db = db
	atomic.StoreUint64(&(a.minBlock), db.GetMinBlock())

	a.status = Running

	return nil
}

// Stop stop the app
func (a *App) Stop() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.status == Stopped {
		return
	}

	a.db.Close()
	a.status = Stopped
}

// Commit will commit something
func (a *App) Commit(data []byte) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if block := UnmarshalHybridBlock(data); block != nil {
		hash := string(crypto.Hash(block.Bytes()))
		if !util.Contain(a.blocks, block.GetNumber()) {
			a.blocks[block.GetNumber()] = block
			a.db.PutBlock(block)
			if block.GetNumber() >= a.getMinBlock() {
				a.hub.Done(hash, nil)
				a.blockCh <- block
			}
			// should parse the hybrid block to different channel block
			if block.GetNumber() == a.getMinBlock() {

			}
		} else {
			a.hub.Done(hash, &event.Result{
				Err: errors.New("Duplicated block"),
			})
		}
	}
}

// Marshal is used in snapshot to get the bytes of data
func (a *App) Marshal() ([]byte, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if len(a.blocks) != 0 {
		return json.Marshal(a.blocks)
	}

	return nil, nil
}

// UnMarshal recover from snapshot
func (a *App) UnMarshal(data []byte) error {
	var blocks map[uint64]*HybridBlock
	if err := json.Unmarshal(data, &blocks); err != nil {
		return err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	for num, block := range blocks {
		if num >= a.getMinBlock() && !util.Contain(a.blocks, num) {
			a.blocks[num] = block
		}
	}

	// Clone the a.blocks to blocks
	data, err := json.Marshal(a.blocks)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &blocks); err != nil {
		return err
	}
	// send the clone blocks
	go func() {
		var nums = make([]uint64, 0)
		for num := range blocks {
			nums = append(nums, num)
		}
		sort.Slice(nums, func(i, j int) bool {
			return nums[i] < nums[j]
		})

		for _, num := range nums {
			block := blocks[num]
			if block.GetNumber() >= a.getMinBlock() {
				a.blockCh <- block
			}
		}

	}()

	return nil
}

func (a *App) watch(block *HybridBlock) error {
	hash := string(crypto.Hash(block.Bytes()))
	res := a.hub.Watch(hash, nil)
	return res.Err
}

// notifyLater provide a mechanism for blockchain system to deal with the block which is too advanced
func (a *App) notifyLater(block *HybridBlock) {
	a.blockCh <- block
}

func (a *App) fetchBlockDone(num uint64) {
	a.lock.Lock()
	defer a.lock.Unlock()

	atomic.StoreUint64(&(a.minBlock), num+1)
	delete(a.blocks, num)
	log.Infof("fetchBlockDone: set minBlock %d + 1", num)
	a.db.SetMinBlock(num + 1)
}

func (a *App) getMinBlock() uint64 {
	return atomic.LoadUint64(&(a.minBlock))
}

// // GetBlock is the implementation of interface
// func (a *App) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
// 	a.lock.Lock()
// 	for i := range a.blocks {
// 		if a.blocks[i].GetNumber() == num {
// 			defer a.lock.Unlock()
// 			log.Infof("consensus/raft/app: get block %d from app.blocks[%s]", num, channelID)
// 			return a.blocks[i], nil
// 		}
// 	}
// 	a.lock.Unlock()
// 	// But block is not in blocks does not mean it is not exist
// 	// todo: can't get block from db
// 	block, _ := a.db.GetBlock(num)
// 	if block != nil {
// 		log.Infof("consensus/raft/app: get block %d from a.db and key is %s", num, channelID)
// 		return block, nil
// 	}

// 	if async {
// 		log.Infof("Watch block %s", fmt.Sprintf("%s:%d", channelID, num))
// 		a.hub.Watch(fmt.Sprintf("%s:%d", channelID, num), nil)
// 		a.lock.Lock()
// 		defer a.lock.Unlock()
// 		for i := range a.blocks {
// 			if a.blocks[i].GetNumber() == num {
// 				log.Infof("consensus/raft/app: get block %d from a.blocks[%s] asynchronously", num, channelID)
// 				return a.blocks[i], nil
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("Block %s:%d is not exist", channelID, num)
// }
