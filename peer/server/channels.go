package server

import (
	"fmt"
	"madledger/common/util"
	"madledger/core"
	"madledger/peer/channel"
	"madledger/peer/config"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"sync"
	"time"
)

// ChannelManager manages all the channels
type ChannelManager struct {
	lock     sync.RWMutex
	db       db.DB
	identity *core.Member

	// signalCh receive stop signal
	signalCh chan bool
	stopCh   chan bool

	// GlobalChannel is the global channel manager
	GlobalChannel *channel.Manager
	// ConfigChannel is the config channel manager
	ConfigChannel *channel.Manager
	// Channels manager all user channels
	Channels map[string]*channel.Manager

	coordinator    *channel.Coordinator
	ordererClients []*orderer.Client
	chainCfg       *config.BlockChainConfig
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager(dbDir string, identity *core.Member, chainCfg *config.BlockChainConfig, ordererClients []*orderer.Client) (*ChannelManager, error) {
	m := new(ChannelManager)
	m.signalCh = make(chan bool, 1)
	m.stopCh = make(chan bool, 1)
	m.Channels = make(map[string]*channel.Manager)
	m.identity = identity
	// set db
	db, err := newDB(dbDir)
	if err != nil {
		return nil, err
	}
	m.db = db
	m.ordererClients = ordererClients
	m.chainCfg = chainCfg
	m.coordinator = channel.NewCoordinator()
	// set global channel manager
	globalManager, err := channel.NewManager(core.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, core.GLOBALCHANNELID), identity, m.db, ordererClients, m.coordinator)
	if err != nil {
		return nil, err
	}
	configManager, err := channel.NewManager(core.CONFIGCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, core.CONFIGCHANNELID), identity, m.db, ordererClients, m.coordinator)
	if err != nil {
		return nil, err
	}
	m.GlobalChannel = globalManager
	m.ConfigChannel = configManager
	for _, channel := range m.db.GetChannels() {
		switch channel {
		case core.GLOBALCHANNELID, core.CONFIGCHANNELID:
		default:
			if !m.hasChannel(channel) {
				m.loadChannel(channel)
			}
		}
	}
	return m, nil
}

// GetTxStatus return the status of tx
func (m *ChannelManager) GetTxStatus(channelID, txID string, async bool) (*db.TxStatus, error) {
	if async {
		return m.db.GetTxStatusAsync(channelID, txID)
	}
	return m.db.GetTxStatus(channelID, txID)
}

// GetTxHistory return all txs of the address
func (m *ChannelManager) GetTxHistory(address []byte) map[string][]string {
	return m.db.GetTxHistory(address)
}

func (m *ChannelManager) start() error {
	updateCh := m.coordinator.RegisterUpdate()
	go m.GlobalChannel.Start()
	go m.ConfigChannel.Start()
	for _, manage := range m.Channels {
		go manage.Start()
	}
	go func() {
		for {
			select {
			case msg := <-updateCh:
				// todo: support channel remove.
				update := msg.(channel.Update)
				if !update.Remove {
					switch update.ID {
					case core.GLOBALCHANNELID, core.CONFIGCHANNELID:
					default:
						if !m.hasChannel(update.ID) {
							manager, err := m.loadChannel(update.ID)
							if err == nil {
								go manager.Start()
							}
						}
					}
				}
			case <-m.signalCh:
				m.stopCh <- true
				return
			}

		}
	}()
	time.Sleep(20 * time.Millisecond)
	return nil
}

// stop will stop all managers
func (m *ChannelManager) stop() {
	log.Info("ChannelManager stop begin")
	m.GlobalChannel.Stop()
	log.Info("GlobalChannel stop")
	m.ConfigChannel.Stop()
	log.Info("ConfigChannel stop")

	m.signalCh <- true
	<-m.stopCh
	for _, manager := range m.Channels {
		manager.Stop()
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
	manager, err := channel.NewManager(channelID, fmt.Sprintf("%s/%s", m.chainCfg.Path, channelID), m.identity, m.db, m.ordererClients, m.coordinator)
	if err != nil {
		return nil, err
	}
	m.Channels[channelID] = manager
	return manager, nil
}
