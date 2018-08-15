package solo

import (
	"errors"
	"fmt"
	"madledger/consensus"
	"madledger/util"
	"sync"
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

func (m *manager) contain(channelID string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return util.Contain(m.channels, channelID)
}

func (m *manager) add(channelID string, cfg consensus.Config) error {
	if m.contain(channelID) {
		return fmt.Errorf("Channel %s is contained aleardy", channelID)
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	channel := newChannel(channelID, cfg, nil)
	m.channels[channelID] = channel
	return nil
}

// todo
func (m *manager) update(channelID string, cfg consensus.Config) error {
	return errors.New("The update is not supported now")
}
