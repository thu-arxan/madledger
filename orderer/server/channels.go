package server

import (
	"fmt"
	cc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	"madledger/consensus"
	"madledger/consensus/solo"
	"madledger/core"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	"madledger/orderer/db"
	"madledger/util"
	"sync"

	"github.com/rs/zerolog/log"
)

// ChannelManager is the manager of channels
type ChannelManager struct {
	chainCfg *config.BlockChainConfig
	db       db.DB
	// Channels manager all user channels
	// maybe can use sync.Map, but the advantage is not significant
	Channels map[string]*channel.Manager
	lock     sync.RWMutex
	// GlobalChannel is the global channel manager
	GlobalChannel *channel.Manager
	// ConfigChannel is the config channel manager
	ConfigChannel *channel.Manager
	Consensus     consensus.Consensus
}

// NewChannelManager is the constructor of ChannelManager
func NewChannelManager(dbDir string, chainCfg *config.BlockChainConfig) (*ChannelManager, error) {
	m := new(ChannelManager)
	m.Channels = make(map[string]*channel.Manager)
	m.chainCfg = chainCfg
	// set db
	db, err := db.NewLevelDB(dbDir)
	if err != nil {
		return nil, err
	}
	m.db = db
	//set config channel manager
	configManager, err := loadConfigChannel(chainCfg.Path, m.db)
	if err != nil {
		return nil, err
	}
	// set global channel manager
	globalManager, err := channel.NewManager(core.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, core.GLOBALCHANNELID), m.db)
	if err != nil {
		return nil, err
	}
	if !globalManager.HasGenesisBlock() {
		log.Info().Msg("Creating genesis block of channel _global")
		// cgb: config channel genesis block
		cgb, err := configManager.GetBlock(0)
		if err != nil {
			return nil, err
		}
		// ggb: global channel genesis block
		ggb, err := gc.CreateGenesisBlock([]*gc.Payload{&gc.Payload{
			ChannelID: core.CONFIGCHANNELID,
			Number:    0,
			Hash:      cgb.Hash(),
		}})
		if err != nil {
			return nil, err
		}
		err = globalManager.AddBlock(ggb)
		if err != nil {
			return nil, err
		}
	}

	m.ConfigChannel = configManager
	m.GlobalChannel = globalManager

	// set consensus
	var channels = make(map[string]consensus.Config, 0)
	cfg := consensus.Config{
		Timeout: 1000,
		MaxSize: 10,
		Number:  0,
		Resume:  false,
	}
	channels[core.GLOBALCHANNELID] = cfg
	channels[core.CONFIGCHANNELID] = cfg
	consensus, err := solo.NewConsensus(channels)
	if err != nil {
		return nil, err
	}
	m.Consensus = consensus

	return m, nil
}

func loadConfigChannel(dir string, db db.DB) (*channel.Manager, error) {
	configManager, err := channel.NewManager(core.CONFIGCHANNELID, fmt.Sprintf("%s/%s", dir, core.CONFIGCHANNELID), db)
	if err != nil {
		return nil, err
	}
	if !configManager.HasGenesisBlock() {
		log.Info().Msg("Creating genesis block of channel _config")
		gb, err := cc.CreateGenesisBlock()
		if err != nil {
			return nil, err
		}
		err = configManager.AddBlock(gb)
		if err != nil {
			return nil, err
		}
	}
	return configManager, nil
}

// TODO
func (manager *ChannelManager) start() error {
	return manager.Consensus.Start()
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
		}
		return nil
	}
}
