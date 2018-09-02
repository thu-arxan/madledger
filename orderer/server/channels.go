package server

import (
	"fmt"
	cc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/solo"
	"madledger/core/types"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	"madledger/orderer/db"
	pb "madledger/protos"
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
	globalManager, err := channel.NewManager(types.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, types.GLOBALCHANNELID), m.db)
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
			ChannelID: types.CONFIGCHANNELID,
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
	channels[types.GLOBALCHANNELID] = cfg
	channels[types.CONFIGCHANNELID] = cfg
	consensus, err := solo.NewConsensus(channels)
	if err != nil {
		return nil, err
	}
	m.Consensus = consensus

	return m, nil
}

// FetchBlock return the block if both channel and block exists
func (manager *ChannelManager) FetchBlock(channelID string, num uint64) (*types.Block, error) {
	cm := manager.getChannelManager(channelID)
	if cm == nil {
		return nil, fmt.Errorf("Channel %s is not exist", channelID)
	}
	return cm.FetchBlock(num)
}

// ListChannels return infos channels
func (manager *ChannelManager) ListChannels(req *pb.ListChannelsRequest) *pb.ChannelInfos {
	infos := new(pb.ChannelInfos)
	if req.System {
		if manager.GlobalChannel != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.GLOBALCHANNELID,
				BlockSize: manager.GlobalChannel.GetBlockSize(),
			})
		}
		if manager.ConfigChannel != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.CONFIGCHANNELID,
				BlockSize: manager.ConfigChannel.GetBlockSize(),
			})
		}
	}
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	for channel := range manager.Channels {
		infos.Channels = append(infos.Channels, &pb.ChannelInfo{
			ChannelID: channel,
		})
	}
	return infos
}

func loadConfigChannel(dir string, db db.DB) (*channel.Manager, error) {
	configManager, err := channel.NewManager(types.CONFIGCHANNELID, fmt.Sprintf("%s/%s", dir, types.CONFIGCHANNELID), db)
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

func (manager *ChannelManager) start() error {
	return manager.Consensus.Start()
}

// stop will stop the consensus
// todo: Not finished yet
func (manager *ChannelManager) stop() {

}

func (manager *ChannelManager) getChannelManager(channelID string) *channel.Manager {
	switch channelID {
	case types.GLOBALCHANNELID:
		return manager.GlobalChannel
	case types.CONFIGCHANNELID:
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
