package server

import (
	"fmt"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/peer/channel"
	"madledger/peer/config"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"sync"
	"time"
)

// ChannelManager manages all the channels
type ChannelManager struct {
	db db.DB
	// Channels manager all user channels
	// maybe can use sync.Map, but the advantage is not significant
	Channels map[string]*channel.Manager
	lock     sync.RWMutex
	// GlobalChannel is the global channel manager
	GlobalChannel *channel.Manager
	// ConfigChannel is the config channel manager
	ConfigChannel *channel.Manager
	ordererClient *orderer.Client
	chainCfg      *config.BlockChainConfig
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager(dbDir string, chainCfg *config.BlockChainConfig, ordererClient *orderer.Client) (*ChannelManager, error) {
	m := new(ChannelManager)
	m.Channels = make(map[string]*channel.Manager)
	// set db
	db, err := db.NewLevelDB(dbDir)
	if err != nil {
		return nil, err
	}
	m.db = db
	m.ordererClient = ordererClient
	m.chainCfg = chainCfg
	// set global channel manager
	globalManager, err := channel.NewManager(types.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, types.GLOBALCHANNELID), m.db, ordererClient)
	if err != nil {
		return nil, err
	}
	configManager, err := channel.NewManager(types.CONFIGCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, types.CONFIGCHANNELID), m.db, ordererClient)
	if err != nil {
		return nil, err
	}
	m.GlobalChannel = globalManager
	m.ConfigChannel = configManager

	return m, nil
}

// GetTxStatus return the status of tx
func (m *ChannelManager) GetTxStatus(channelID, txID string) (*db.TxStatus, error) {
	return m.db.GetTxStatus(channelID, txID)
}

func (m *ChannelManager) start() error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	go m.GlobalChannel.Start()
	go m.ConfigChannel.Start()
	for {
		select {
		case <-ticker.C:
			channels, err := m.ordererClient.ListChannels()
			if err == nil {
				for _, channel := range channels {
					if !m.hasChannel(channel) {
						manager, err := m.loadChannel(channel)
						if err == nil {
							go manager.Start()
						}
					}
				}
			}
		}
	}
}

// hasChannel return if a channel exist
func (m *ChannelManager) hasChannel(channelID string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if util.Contain(m.Channels, channelID) {
		return true
	}
	return false
}

// loadChannel load a channel
func (m *ChannelManager) loadChannel(channelID string) (*channel.Manager, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if util.Contain(m.Channels, channelID) {
		return m.Channels[channelID], nil
	}
	manager, err := channel.NewManager(channelID, fmt.Sprintf("%s/%s", m.chainCfg.Path, channelID), m.db, m.ordererClient)
	if err != nil {
		return nil, err
	}
	m.Channels[channelID] = manager
	return manager, nil
}
