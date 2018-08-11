package server

import (
	"madledger/orderer/channel"
	"madledger/orderer/db"
)

// ChannelManager is the manager of channels
type ChannelManager struct {
	db *db.DB
	// Channels manager all user channels
	Channels map[string]*channel.Manager
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
