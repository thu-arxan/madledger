package server

import (
	"fmt"
	"madledger/core"
	"madledger/orderer/channel"
	"madledger/orderer/db"
	"madledger/util"
	"sync"
)

// ChannelManager is the manager of channels
type ChannelManager struct {
	db *db.DB
	// Channels manager all user channels
	// maybe can use sync.Map, but the advantage is not significant
	Channels map[string]*channel.Manager
	lock     sync.RWMutex
	// GlobalChannel is the global channel manager
	GlobalChannel *channel.Manager
	// ConfigChannel is the config channel manager
	ConfigChannel *channel.Manager
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager(dir string) (*ChannelManager, error) {
	m := new(ChannelManager)
	m.Channels = make(map[string]*channel.Manager)
	// set db
	db, err := db.NewGolevelDB(dir)
	if err != nil {
		return nil, err
	}
	m.db = &db
	// set global channel manager
	globalManager, err := channel.NewManager("_global", m.db)
	if err != nil {
		return nil, err
	}
	m.GlobalChannel = globalManager
	//set config channel manager
	configManager, err := channel.NewManager("_config", m.db)
	if err != nil {
		return nil, err
	}
	m.ConfigChannel = configManager
	return m, nil
}

// FetchBlock return the block if both channel and block exists
func (manager *ChannelManager) FetchBlock(channelID string, num uint64) (*core.Block, error) {
	cm := manager.getChannelManager(channelID)
	if cm == nil {
		return nil, fmt.Errorf("Channel %s is not exist", channelID)
	}
	return cm.FetchBlock(num)
}

func (manager *ChannelManager) getChannelManager(channelID string) *channel.Manager {
	switch channelID {
	case "_global":
		return manager.GlobalChannel
	case "_config":
		return manager.ConfigChannel
	default:
		manager.lock.RLock()
		defer manager.lock.RUnlock()
		if util.Contain(manager.Channels, channelID) {
			return manager.Channels[channelID]
		} else {
			return nil
		}
	}
}
