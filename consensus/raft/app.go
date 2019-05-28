package raft

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	core "madledger/core/types"
	"sort"
	"sync"
	"sync/atomic"
)

// App is the application
type App struct {
	lock   sync.Mutex
	cfg    *Config
	status int32 // only Running and Stopped

	blocks  map[uint64]*core.Block
	blockCh chan *core.Block
	hub     *event.Hub
	// minBlock is the min block number that the blockchain system needed
	minBlock uint64
	db       *DB
}

// NewApp is the constructor of App
func NewApp(cfg *Config) (*App, error) {
	return &App{
		cfg:     cfg,
		blocks:  make(map[uint64]*core.Block),
		blockCh: make(chan *core.Block, 2048),
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

	var block core.Block
	if err := json.Unmarshal(data, &block); err == nil {
		hash := string(crypto.Hash(block.Bytes()))
		if !util.Contain(a.blocks, block.GetNumber()) {
			a.blocks[block.GetNumber()] = &block
			if block.GetNumber() >= a.getMinBlock() {
				a.db.AddBlock(&block)
				a.hub.Done(hash, nil)
				a.blockCh <- &block
				// a.sendBlocks()
			}
		} else {
			// todo: need finish this
			// a.hub.Done(hash, fmt.Errorf("[%d] Duplicated block", a.cfg.id))
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
	var blocks map[uint64]*core.Block
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
		// fmt.Println(a.cfg.id, ":", nums)
		for _, num := range nums {
			block := blocks[num]
			if block.GetNumber() >= a.getMinBlock() {
				a.blockCh <- block
			}
		}

	}()

	return nil
}

func (a *App) watch(block *core.Block) error {
	hash := string(crypto.Hash(block.Bytes()))
	res := a.hub.Watch(hash, nil)
	return res.Err
}

// notifyLater provide a mechanism for blockchain system to deal with the block which is too advanced
func (a *App) notifyLater(block *core.Block) {
	a.blockCh <- block
}

func (a *App) fetchBlockDone(num uint64) {
	a.lock.Lock()
	defer a.lock.Unlock()

	atomic.StoreUint64(&(a.minBlock), num+1)
	delete(a.blocks, num)
	a.db.SetMinBlock(num + 1)
}

func (a *App) getMinBlock() uint64 {
	return atomic.LoadUint64(&(a.minBlock))
}

// sendBlocks may be a good way, but it need to be designed better
// todo: redesign the function, maybe need change notifyLater function
// Note: The function is not used now.
func (a *App) sendBlocks() {
	// a.lock.Lock()
	// defer a.lock.Unlock()
	var nums = make([]uint64, 0)
	for num := range a.blocks {
		nums = append(nums, num)
	}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	minBlock := a.getMinBlock()
	if util.Contain(a.blocks, minBlock) {
		// fmt.Println(a.cfg.id, ":", nums)
		for _, num := range nums {
			block := a.blocks[num]
			if block.GetNumber() >= minBlock {
				a.blockCh <- block
			}
		}
	} else {
		log.Infof("[%d] Miss block %d:\n", a.cfg.id, minBlock)
		log.Info(nums)
	}
}
