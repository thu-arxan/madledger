package server

import "madledger/orderer/channel"

// ChannelManager is the manager of channels
type ChannelManager struct {
	Channels map[string]*channel.Manager
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager() (*ChannelManager, error) {
	m := new(ChannelManager)
	m.Channels = make(map[string]*channel.Manager)
	return m, nil
}
