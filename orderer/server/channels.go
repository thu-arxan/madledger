package server

import (
	"encoding/json"
	"fmt"
	cc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/solo"
	"madledger/core/types"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	"madledger/orderer/db"
	pb "madledger/protos"
	"strings"
	"sync"
	"time"
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
		log.Info("Creating genesis block of channel _global")
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

	// then load user channels
	userChannels, err := loadUserChannels(chainCfg.Path, db)
	if err != nil {
		return nil, err
	}
	for channelID, manager := range userChannels {
		m.Channels[channelID] = manager
	}
	// set consensus
	var channels = make(map[string]consensus.Config, 0)
	cfg := consensus.Config{
		Timeout: 100,
		MaxSize: 10,
		Number:  1,
		Resume:  false,
	}
	channels[types.GLOBALCHANNELID] = cfg
	channels[types.CONFIGCHANNELID] = cfg
	// set consensus of user channels
	for channelID := range userChannels {
		channels[channelID] = consensus.Config{
			Timeout: 1000,
			MaxSize: 10,
			Number:  1,
			Resume:  false,
		}
	}
	consensus, err := solo.NewConsensus(channels)
	if err != nil {
		return nil, err
	}
	m.Consensus = consensus

	return m, nil
}

// FetchBlock return the block if both channel and block exists
func (manager *ChannelManager) FetchBlock(channelID string, num uint64, async bool) (*types.Block, error) {
	cm, err := manager.getChannelManager(channelID)
	if err != nil {
		return nil, err
	}
	if async {
		return cm.FetchBlockAsync(num)
	}
	return cm.FetchBlock(num)
}

// ListChannels return infos channels
func (manager *ChannelManager) ListChannels(req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	pk, err := crypto.NewPublicKey(req.PK)
	if err != nil {
		return &pb.ChannelInfos{}, err
	}
	member, err := types.NewMember(pk, "")
	if err != nil {
		return &pb.ChannelInfos{}, err
	}
	infos := new(pb.ChannelInfos)
	if req.System {
		if manager.GlobalChannel != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.GLOBALCHANNELID,
				BlockSize: manager.GlobalChannel.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
		if manager.ConfigChannel != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.CONFIGCHANNELID,
				BlockSize: manager.ConfigChannel.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
	}
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	for channel, channelManager := range manager.Channels {
		if channelManager.IsMember(member) {
			identity := pb.Identity_MEMBER
			if channelManager.IsAdmin(member) {
				identity = pb.Identity_ADMIN
			}
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: channel,
				BlockSize: channelManager.GetBlockSize(),
				Identity:  identity,
			})
		}
	}
	return infos, nil
}

// CreateChannel try to create a channel
func (manager *ChannelManager) CreateChannel(tx *types.Tx) (*pb.ChannelInfo, error) {
	err := manager.createChannel(tx)
	if err != nil {
		return nil, err
	}
	return &pb.ChannelInfo{}, nil
}

// createChannel try to create a channel
// However, this should check if the channel exist and should be thread safety.
// todo: First add a tx then create channel
func (manager *ChannelManager) createChannel(tx *types.Tx) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return err
	}
	var channelID = payload.ChannelID
	switch channelID {
	case types.GLOBALCHANNELID:
	case types.CONFIGCHANNELID:
		return fmt.Errorf("Channel %s is aleardy exist", channelID)
	default:
		if !isLegalChannelName(channelID) {
			return fmt.Errorf("%s is not a legal channel name", channelID)
		}
		if util.Contain(manager.Channels, channelID) {
			return fmt.Errorf("Channel %s is aleardy exist", channelID)
		}
	}
	// then try to create a channel
	_, err = channel.NewManager(channelID, fmt.Sprintf("%s/%s", manager.chainCfg.Path, channelID), manager.db)
	if err != nil {
		return err
	}
	// then send a tx to config channel
	// But the manager should not AddTx by consensus, because the confirm
	// of consensus is not the final confirm.
	err = manager.ConfigChannel.AddTx(tx)
	if err != nil {
		return err
	}

	// then start the consensus
	err = manager.Consensus.AddChannel(channelID, consensus.Config{
		Timeout: 1000,
		MaxSize: 10,
		Number:  0,
		Resume:  false,
	})
	channel, err := channel.NewManager(channelID, fmt.Sprintf("%s/%s", manager.chainCfg.Path, channelID), manager.db)
	if err != nil {
		return err
	}
	// create genesis block here
	// The genesis only contain the create tx now.
	genesisBlock := types.NewBlock(channelID, 0, types.GenesisBlockPrevHash, []*types.Tx{tx})
	err = channel.AddBlock(genesisBlock)
	if err != nil {
		return err
	}
	// then start the channel
	go func() {
		channel.Start(manager.Consensus, manager.GlobalChannel)
	}()
	manager.Channels[channelID] = channel
	return err
}

// AddTx add a tx
func (manager *ChannelManager) AddTx(tx *types.Tx) error {
	channel, err := manager.getChannelManager(tx.Data.ChannelID)
	if err != nil {
		return err
	}
	return channel.AddTx(tx)
}

// loadConfigChannel load the config channel("_config")
func loadConfigChannel(dir string, db db.DB) (*channel.Manager, error) {
	configManager, err := channel.NewManager(types.CONFIGCHANNELID, fmt.Sprintf("%s/%s", dir, types.CONFIGCHANNELID), db)
	if err != nil {
		return nil, err
	}
	if !configManager.HasGenesisBlock() {
		log.Info("Creating genesis block of channel _config")
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

func loadUserChannels(dir string, db db.DB) (map[string]*channel.Manager, error) {
	var managers = make(map[string]*channel.Manager)
	channels := db.ListChannel()
	for _, channelID := range channels {
		if !strings.HasPrefix(channelID, "_") {
			manager, err := channel.NewManager(channelID, fmt.Sprintf("%s/%s", dir, channelID), db)
			if err != nil {
				return nil, err
			}
			managers[channelID] = manager
		}
	}
	return managers, nil
}

// start the manager
func (manager *ChannelManager) start() error {
	// start consensus
	err := manager.Consensus.Start()
	if err != nil {
		return err
	}

	go manager.GlobalChannel.Start(manager.Consensus, nil)
	go manager.ConfigChannel.Start(manager.Consensus, manager.GlobalChannel)
	// also, start others channels
	for _, channelManager := range manager.Channels {
		go channelManager.Start(manager.Consensus, manager.GlobalChannel)
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

// stop will stop the consensus
func (manager *ChannelManager) stop() error {
	defer manager.db.Close()
	return manager.Consensus.Stop()
}

// getChannelManager return the manager of the channel
func (manager *ChannelManager) getChannelManager(channelID string) (*channel.Manager, error) {
	switch channelID {
	case types.GLOBALCHANNELID:
		return manager.GlobalChannel, nil
	case types.CONFIGCHANNELID:
		return manager.ConfigChannel, nil
	default:
		manager.lock.RLock()
		defer manager.lock.RUnlock()
		if util.Contain(manager.Channels, channelID) {
			return manager.Channels[channelID], nil
		}
	}
	return nil, fmt.Errorf("Channel %s is not exist", channelID)
}
