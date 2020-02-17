package eraft

import (
	"encoding/json"
	"errors"
	"sync"
)

// App is the application
type App struct {
	lock sync.RWMutex

	cfg    *EraftConfig
	status int32 // only Running or Stopped

	channelsLock sync.RWMutex
	channels     map[string]*channel // channelID => channel storage

	db *DB
}

// NewApp is the constructor of App
func NewApp(cfg *EraftConfig) (*App, error) {
	return &App{
		cfg:      cfg,
		channels: make(map[string]*channel),
		status:   Stopped,
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
	block := UnmarshalBlock(data)
	if block == nil {
		return
	}
	channel := a.getChannel(block.ChannelID)
	channel.addBlock(block)
}

// Marshal is used in snapshot to get the bytes of data
func (a *App) Marshal() ([]byte, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.marshalChannelBlocks()
}

// UnMarshal recover from snapshot
// todo: is it necessary to do this
func (a *App) UnMarshal(data []byte) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.unmarshalChannelBlocks(data)
}

func (a *App) watch(block *Block) error {
	channel := a.getChannel(block.ChannelID)
	return channel.watch(block)
}

func (a *App) blockCh(channelID string) chan *Block {
	channel := a.getChannel(channelID)
	return channel.BlockCh()
}

// notifyLater provide a mechanism for blockchain system to deal with the block which is too advanced
func (a *App) notifyLater(block *Block) {
	channel := a.getChannel(block.ChannelID)
	channel.notifyLater(block)
}

func (a *App) fetchBlockDone(channelID string, num uint64) {
	channel := a.getChannel(channelID)
	channel.fetchBlockDone(num)
}

func (a *App) getChannel(channelID string) *channel {
	a.channelsLock.RLock()
	channel := a.channels[channelID]
	a.channelsLock.RUnlock()

	if channel == nil {
		channel = newChannel(a.cfg.id, channelID, a.db)
		a.channelsLock.Lock()
		a.channels[channelID] = channel
		a.channelsLock.Unlock()
	}
	return channel
}

func (a *App) marshalChannelBlocks() ([]byte, error) {
	a.channelsLock.RLock()
	defer a.channelsLock.RUnlock()

	var err error

	channelBlocks := make(map[string][]byte)

	for channelID, channel := range a.channels {
		channelBlocks[channelID], err = channel.marshal()
		if err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(channelBlocks)
	return data, err
}

func (a *App) unmarshalChannelBlocks(data []byte) error {
	channelBlocks := make(map[string][]byte)
	if err := json.Unmarshal(data, &channelBlocks); err != nil {
		return err
	}

	for channelID, blockData := range channelBlocks {
		channel := a.getChannel(channelID)
		if err := channel.unmarshal(blockData); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) setChainNum(channelID string, num uint64) {
	a.db.SetChainNum(channelID, num)
}

func (a *App) getChainNum(channelID string) uint64 {
	return a.db.GetChainNum(channelID)
}
