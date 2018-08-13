package server

import (
	"fmt"
	"madledger/core"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	"madledger/orderer/db"
	"madledger/util"
	"sync"
)

// ChannelManager is the manager of channels
type ChannelManager struct {
	chainCfg *config.BlockChainConfig
	db       *db.DB
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
func NewChannelManager(dbDir string, chainCfg *config.BlockChainConfig) (*ChannelManager, error) {
	m := new(ChannelManager)
	m.Channels = make(map[string]*channel.Manager)
	m.chainCfg = chainCfg
	// set db
	db, err := db.NewGolevelDB(dbDir)
	if err != nil {
		return nil, err
	}
	m.db = &db
	// set global channel manager
	globalManager, err := channel.NewManager(core.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, core.GLOBALCHANNELID), m.db)
	if err != nil {
		return nil, err
	}
	m.GlobalChannel = globalManager
	//set config channel manager
	configManager, err := channel.NewManager(core.CONFIGCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, core.CONFIGCHANNELID), m.db)
	if err != nil {
		return nil, err
	}
	m.ConfigChannel = configManager
	return m, nil
}

// TODO
func (manager *ChannelManager) start() error {
	return nil
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
	case core.GLOBALCHANNELID:
		return manager.GlobalChannel
	case core.CONFIGCHANNELID:
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
