// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package solo

import (
	"errors"
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"sync"
	"time"
)

type manager struct {
	lock     sync.RWMutex
	channels map[string]*channel
}

func newManager() *manager {
	m := new(manager)
	m.channels = make(map[string]*channel, 0)
	return m
}

func (m *manager) start() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, channel := range m.channels {
		go channel.start()
	}
	return nil
}

func (m *manager) stop() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var wg sync.WaitGroup

	for _, c := range m.channels {
		wg.Add(1)
		go func(c *channel) {
			defer wg.Done()
			c.Stop()
		}(c)
	}
	wg.Wait()
	return nil
}

func (m *manager) contain(channelID string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return util.Contain(m.channels, channelID)
}

func (m *manager) get(channelID string) (*channel, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if util.Contain(m.channels, channelID) {
		return m.channels[channelID], nil
	}
	return nil, fmt.Errorf("The channel %s is not exist", channelID)
}

// AddTx add a tx
func (m *manager) AddTx(tx *core.Tx) error {
	channel, err := m.get(tx.Data.ChannelID)
	if err != nil {
		return err
	}
	return channel.AddTx(tx)
}

func (m *manager) add(channelID string, cfg consensus.Config) error {
	if m.contain(channelID) {
		return fmt.Errorf("Channel %s is contained aleardy", channelID)
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	channel := newChannel(channelID, cfg)
	m.channels[channelID] = channel
	return nil
}

func (m *manager) startChannel(channelID string) error {
	channel, err := m.get(channelID)
	if err != nil {
		return err
	}
	if channel.initialized() {
		return fmt.Errorf("Channel %s is starting aleardy", channelID)
	}
	go channel.start()
	time.Sleep(20 * time.Millisecond)
	return nil
}

// todo: update consensus config is not finished yet
func (m *manager) update(channelID string, cfg consensus.Config) error {
	return errors.New("The update is not supported now")
}
