package channel

import (
	"encoding/json"
	"fmt"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/solo"
	"madledger/consensus/tendermint"
	"madledger/core/types"
	"madledger/orderer/config"
	"madledger/orderer/db"
	"strings"
	"sync"
	"time"

	cc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	ct "madledger/consensus/tendermint"
	pb "madledger/protos"
)

// TODO: These codes is too complex and may contains lots of bugs. So conside rewrite them right away.

// Coordinator responsible for coordination of managers
type Coordinator struct {
	chainCfg *config.BlockChainConfig
	db       db.DB

	lock sync.RWMutex
	// Channels manager all user channels
	// maybe can use sync.Map, but the advantage is not significant
	Managers map[string]*Manager
	// GM is the global channel manager
	GM *Manager
	// CM is the config channel manager
	CM        *Manager
	Consensus consensus.Consensus
}

// NewCoordinator is the constructor of Coordinator
func NewCoordinator(dbDir string, chainCfg *config.BlockChainConfig, consensusCfg *config.ConsensusConfig) (*Coordinator, error) {
	c := new(Coordinator)
	c.Managers = make(map[string]*Manager)
	c.chainCfg = chainCfg
	// set db
	db, err := db.NewLevelDB(dbDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to load db at %s because %s", dbDir, err.Error())
	}
	c.db = db
	//set config channel manager
	configManager, err := c.loadConfigChannel(chainCfg.Path)
	if err != nil {
		return nil, err
	}
	// set global channel manager
	globalManager, err := NewManager(types.GLOBALCHANNELID, fmt.Sprintf("%s/%s", chainCfg.Path, types.GLOBALCHANNELID), c)
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

	c.CM = configManager
	c.GM = globalManager

	// then load user channels
	userChannels, err := c.loadUserChannels(chainCfg.Path)
	if err != nil {
		return nil, err
	}
	for channelID, manager := range userChannels {
		c.Managers[channelID] = manager
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
	switch consensusCfg.Type {
	case config.SOLO:
		consensus, err := solo.NewConsensus(channels)
		if err != nil {
			return nil, err
		}
		c.Consensus = consensus
	case config.BFT:
		// TODO: Not finished yet
		consensus, err := tendermint.NewConsensus(channels, &ct.Config{
			Port: ct.Port{
				P2P: consensusCfg.BFT.Port.P2P,
				RPC: consensusCfg.BFT.Port.RPC,
				App: consensusCfg.BFT.Port.APP,
			},
			Dir:        consensusCfg.BFT.Path,
			P2PAddress: consensusCfg.BFT.P2PAddress,
		})
		if err != nil {
			return nil, err
		}
		c.Consensus = consensus
	default:
		return nil, fmt.Errorf("Unsupport consensus type:%d", consensusCfg.Type)
	}
	return c, nil
}

// Start the coordinator
func (c *Coordinator) Start() error {
	// start consensus
	err := c.Consensus.Start()
	if err != nil {
		return err
	}

	go c.GM.Start()
	go c.CM.Start()
	// also, start others channels
	for _, channelManager := range c.Managers {
		go channelManager.Start()
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

// Stop will stop the consensus
func (c *Coordinator) Stop() error {
	defer c.db.Close()
	return c.Consensus.Stop()
}

// FetchBlock return the block if both channel and block exists
func (c *Coordinator) FetchBlock(channelID string, num uint64, async bool) (*types.Block, error) {
	cm, err := c.getChannelManager(channelID)
	if err != nil {
		return nil, err
	}
	if async {
		return cm.FetchBlockAsync(num)
	}
	return cm.FetchBlock(num)
}

// ListChannels return infos channels
func (c *Coordinator) ListChannels(req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
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
		if c.GM != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.GLOBALCHANNELID,
				BlockSize: c.GM.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
		if c.CM != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: types.CONFIGCHANNELID,
				BlockSize: c.CM.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	for channel, channelManager := range c.Managers {
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
func (c *Coordinator) CreateChannel(tx *types.Tx) (*pb.ChannelInfo, error) {
	err := c.createChannel(tx)
	if err != nil {
		return nil, err
	}
	return &pb.ChannelInfo{}, nil
}

// AddTx add a tx
func (c *Coordinator) AddTx(tx *types.Tx) error {
	channel, err := c.getChannelManager(tx.Data.ChannelID)
	if err != nil {
		return err
	}
	return channel.AddTx(tx)
}

// createChannel try to create a channel
// However, this should check if the channel exist and should be thread safety.
// todo: First add a tx then create channel
func (c *Coordinator) createChannel(tx *types.Tx) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var payload cc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return err
	}
	var channelID = payload.ChannelID
	switch channelID {
	case types.GLOBALCHANNELID:
		return fmt.Errorf("Channel %s is aleardy exist", channelID)
	case types.CONFIGCHANNELID:
		return fmt.Errorf("Channel %s is aleardy exist", channelID)
	default:
		if !isLegalChannelName(channelID) {
			return fmt.Errorf("%s is not a legal channel name", channelID)
		}
		if util.Contain(c.Managers, channelID) {
			return fmt.Errorf("Channel %s is aleardy exist", channelID)
		}
	}

	// then send a tx to config channel
	// But the manager should not AddTx by consensus, because the confirm
	// of consensus is not the final confirm.
	err = c.CM.AddTx(tx)
	if err != nil {
		return err
	}

	// then start the consensus
	err = c.Consensus.AddChannel(channelID, consensus.Config{
		Timeout: 1000,
		MaxSize: 10,
		Number:  0,
		Resume:  false,
	})
	channel, err := NewManager(channelID, fmt.Sprintf("%s/%s", c.chainCfg.Path, channelID), c)
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
		channel.Start()
	}()
	c.Managers[channelID] = channel
	return err
}

// getChannelManager return the manager of the channel
func (c *Coordinator) getChannelManager(channelID string) (*Manager, error) {
	switch channelID {
	case types.GLOBALCHANNELID:
		return c.GM, nil
	case types.CONFIGCHANNELID:
		return c.CM, nil
	default:
		c.lock.RLock()
		defer c.lock.RUnlock()
		if util.Contain(c.Managers, channelID) {
			return c.Managers[channelID], nil
		}
	}
	return nil, fmt.Errorf("Channel %s is not exist", channelID)
}

// loadConfigChannel load the config channel("_config")
func (c *Coordinator) loadConfigChannel(dir string) (*Manager, error) {
	configManager, err := NewManager(types.CONFIGCHANNELID, fmt.Sprintf("%s/%s", dir, types.CONFIGCHANNELID), c)
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

func (c *Coordinator) loadUserChannels(dir string) (map[string]*Manager, error) {
	var managers = make(map[string]*Manager)
	channels := c.db.ListChannel()
	for _, channelID := range channels {
		if !strings.HasPrefix(channelID, "_") {
			manager, err := NewManager(channelID, fmt.Sprintf("%s/%s", dir, channelID), c)
			if err != nil {
				return nil, err
			}
			managers[channelID] = manager
		}
	}
	return managers, nil
}
