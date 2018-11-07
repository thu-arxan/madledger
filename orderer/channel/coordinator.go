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
	CM *Manager

	Consensus consensus.Consensus
}

// NewCoordinator is the constructor of Coordinator
func NewCoordinator(dbDir string, chainCfg *config.BlockChainConfig, consensusCfg *config.ConsensusConfig) (*Coordinator, error) {
	var err error

	c := new(Coordinator)
	c.Managers = make(map[string]*Manager)
	c.chainCfg = chainCfg
	// set db
	c.db, err = db.NewLevelDB(dbDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to load db at %s because %s", dbDir, err.Error())
	}
	//set system channels like config and global
	err = c.loadSystemChannel()
	if err != nil {
		return nil, err
	}
	// then load user channels
	err = c.loadUserChannel()
	if err != nil {
		return nil, err
	}
	// set consensus
	err = c.setConsensus(consensusCfg)
	if err != nil {
		return nil, err
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

	err = c.CM.AddTx(tx)
	if err != nil {
		return err
	}
	c.db.WatchChannel(channelID)

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

// loadSystemChannel will load config channel and global channel
func (c *Coordinator) loadSystemChannel() error {
	if err := c.loadConfigChannel(); err != nil {
		return err
	}

	if err := c.loadGlobalChannel(); err != nil {
		return err
	}
	return nil
}

// loadConfigChannel load the config channel("_config")
func (c *Coordinator) loadConfigChannel() error {
	var err error
	c.CM, err = NewManager(types.CONFIGCHANNELID, c)
	if err != nil {
		return err
	}
	if !c.CM.HasGenesisBlock() {
		log.Info("Creating genesis block of channel _config")
		gb, err := cc.CreateGenesisBlock()
		if err != nil {
			return err
		}
		err = c.CM.AddBlock(gb)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadGlobalChannel load the global channel("_global")
// Note: loadGlobalChannel must call after loadConfigChannel
func (c *Coordinator) loadGlobalChannel() error {
	var err error
	c.GM, err = NewManager(types.GLOBALCHANNELID, c)
	if err != nil {
		return err
	}
	if !c.GM.HasGenesisBlock() {
		log.Info("Creating genesis block of channel _global")
		// cgb: config channel genesis block
		cgb, err := c.CM.GetBlock(0)
		if err != nil {
			return err
		}
		// ggb: global channel genesis block
		ggb, err := gc.CreateGenesisBlock([]*gc.Payload{&gc.Payload{
			ChannelID: types.CONFIGCHANNELID,
			Number:    0,
			Hash:      cgb.Hash(),
		}})
		if err != nil {
			return err
		}
		err = c.GM.AddBlock(ggb)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadUserChannel load all user channels
func (c *Coordinator) loadUserChannel() error {
	channels := c.db.ListChannel()
	for _, channelID := range channels {
		if !strings.HasPrefix(channelID, "_") {
			manager, err := NewManager(channelID, c)
			if err != nil {
				return err
			}
			c.Managers[channelID] = manager
		}
	}
	return nil
}

// setConsensus set consensus according to the consensus config
func (c *Coordinator) setConsensus(cfg *config.ConsensusConfig) error {
	// set consensus
	var channels = make(map[string]consensus.Config, 0)
	defaultCfg := consensus.Config{
		Timeout: 100,
		MaxSize: 10,
		Number:  1,
		Resume:  false,
	}
	channels[types.GLOBALCHANNELID] = defaultCfg
	channels[types.CONFIGCHANNELID] = defaultCfg
	// set consensus of user channels
	for channelID := range c.Managers {
		channels[channelID] = defaultCfg
	}
	switch cfg.Type {
	case config.SOLO:
		consensus, err := solo.NewConsensus(channels)
		if err != nil {
			return err
		}
		c.Consensus = consensus
	case config.BFT:
		// TODO: Not finished yet
		consensus, err := tendermint.NewConsensus(channels, &ct.Config{
			Port: ct.Port{
				P2P: cfg.BFT.Port.P2P,
				RPC: cfg.BFT.Port.RPC,
				App: cfg.BFT.Port.APP,
			},
			Dir:        cfg.BFT.Path,
			P2PAddress: cfg.BFT.P2PAddress,
		})
		if err != nil {
			return err
		}
		c.Consensus = consensus
	default:
		return fmt.Errorf("Unsupport consensus type:%d", cfg.Type)
	}
	return nil
}
