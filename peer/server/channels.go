// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
	path     string

	// signalCh receive stop signal
	signalCh chan bool
	stopCh   chan bool

	// GlobalChannel is the global channel manager
	GlobalChannel *channel.Manager
	// ConfigChannel is the config channel manager
	ConfigChannel *channel.Manager
	// AssetChannel is the asset channel manager
	AssetChannel   *channel.Manager
	// Channels manager all user channels
	Channels map[string]*channel.Manager

	coordinator    *channel.Coordinator
	ordererClients []*orderer.Client
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager(cfg *config.Config) (*ChannelManager, error) {
	m := new(ChannelManager)
	var err error
	m.signalCh = make(chan bool, 1)
	m.stopCh = make(chan bool, 1)
	m.Channels = make(map[string]*channel.Manager)
	// set identity
	if m.identity, err = cfg.GetIdentity(); err != nil {
		return nil, err
	}
	// set db
	db, err := newDB(cfg.DB.LevelDB.Dir)
	if err != nil {
		return nil, err
	}
	m.db = db
	// set path
	m.path = cfg.BlockChain.Path
	// set order clients
	if m.ordererClients, err = getOrdererClients(cfg); err != nil {
		return nil, err
	}
	m.coordinator = channel.NewCoordinator()
	if err := m.loadChannels(); err != nil {
		return nil, err
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
	go m.AssetChannel.Start()
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
					case core.GLOBALCHANNELID, core.CONFIGCHANNELID, core.ASSETCHANNELID:
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
	m.AssetChannel.Stop()
	log.Info("AccountChannel stop")

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

// load system channels and user channels
func (m *ChannelManager) loadChannels() error {
	// set global channel manager
	globalManager, err := channel.NewManager(core.GLOBALCHANNELID, fmt.Sprintf("%s/%s", m.path, core.GLOBALCHANNELID), m.identity, m.db, m.ordererClients, m.coordinator)
	if err != nil {
		return err
	}
	configManager, err := channel.NewManager(core.CONFIGCHANNELID, fmt.Sprintf("%s/%s", m.path, core.CONFIGCHANNELID), m.identity, m.db, m.ordererClients, m.coordinator)
	if err != nil {
		return err
	}
	assetManager, err := channel.NewManager(core.ASSETCHANNELID, fmt.Sprintf("%s/%s", m.path, core.ASSETCHANNELID), m.identity, m.db, m.ordererClients, m.coordinator)
	if err != nil {
		return err
	}
	m.GlobalChannel = globalManager
	m.ConfigChannel = configManager
	m.AssetChannel = assetManager
	for _, channel := range m.db.GetChannels() {
		switch channel {
		case core.GLOBALCHANNELID, core.CONFIGCHANNELID, core.ASSETCHANNELID:
		default:
			m.loadChannel(channel)
		}
	}
	return nil
}

// loadChannel load a channel
func (m *ChannelManager) loadChannel(channelID string) (*channel.Manager, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if util.Contain(m.Channels, channelID) {
		return m.Channels[channelID], nil
	}
	manager, err := channel.NewManager(channelID, fmt.Sprintf("%s/%s", m.path, channelID), m.identity, m.db, m.ordererClients, m.coordinator)
	if err != nil {
		return nil, err
	}
	m.Channels[channelID] = manager
	return manager, nil
}
