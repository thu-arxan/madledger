package channel

import (
	"encoding/json"
	"fmt"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/consensus"
	raft "madledger/consensus/raft"
	"madledger/consensus/solo"
	"madledger/consensus/tendermint"
	"madledger/core"
	"madledger/orderer/config"
	"madledger/orderer/db"
	"strings"
	"sync"
	"time"

	ac "madledger/blockchain/asset"
	bc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	ct "madledger/consensus/tendermint"
	pb "madledger/protos"
)

// Coordinator responsible for coordination of managers
type Coordinator struct {
	chainCfg *config.BlockChainConfig
	db       db.DB

	lock sync.RWMutex
	// Channels manager all user channels
	managerLock sync.RWMutex
	Managers    map[string]*Manager
	// GM is the global channel manager
	GM *Manager
	// CM is the config channel manager
	CM *Manager
	// AM is the asset channel manager
	AM *Manager

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
	time.Sleep(100 * time.Millisecond)
	go c.CM.Start()
	time.Sleep(100 * time.Millisecond)
	go c.AM.Start()
	time.Sleep(100 * time.Millisecond)
	for _, channelManager := range c.Managers {
		// 开启manager之前，应该判断manager的init是否为true
		if !channelManager.init {
			log.Infof("coordinator/Start: start channel %s", channelManager.ID)
			go channelManager.Start()
		}
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
func (c *Coordinator) FetchBlock(channelID string, num uint64, async bool) (*core.Block, error) {
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
	member, err := core.NewMember(pk, "")
	if err != nil {
		return &pb.ChannelInfos{}, err
	}
	infos := new(pb.ChannelInfos)
	if req.System {
		if c.GM != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: core.GLOBALCHANNELID,
				BlockSize: c.GM.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
		if c.CM != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: core.CONFIGCHANNELID,
				BlockSize: c.CM.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
		if c.AM != nil {
			infos.Channels = append(infos.Channels, &pb.ChannelInfo{
				ChannelID: core.ASSETCHANNELID,
				BlockSize: c.AM.GetBlockSize(),
				Identity:  pb.Identity_MEMBER,
			})
		}
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	c.managerLock.RLock()
	defer c.managerLock.RUnlock()

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
func (c *Coordinator) CreateChannel(tx *core.Tx) (*pb.ChannelInfo, error) {
	err := c.createChannel(tx)
	if err != nil {
		return nil, err
	}
	return &pb.ChannelInfo{}, nil
}

// AddTx add a tx
func (c *Coordinator) AddTx(tx *core.Tx) error {
	channel, err := c.getChannelManager(tx.Data.ChannelID)
	if err != nil {
		return err
	}
	return channel.AddTx(tx)
}

func (c *Coordinator) setChannel(channelID string, manager *Manager) {
	c.managerLock.Lock()
	defer c.managerLock.Unlock()
	c.Managers[channelID] = manager
}

// createChannel try to create a channel
// However, this should check if the channel exist and should be thread safety.
func (c *Coordinator) createChannel(tx *core.Tx) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var payload bc.Payload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return err
	}

	var channelID = payload.ChannelID
	log.Infof("Create channel %s", channelID)
	switch channelID {
	case core.GLOBALCHANNELID:
		return fmt.Errorf("Channel %s is already exist", channelID)
	case core.CONFIGCHANNELID:
		return fmt.Errorf("Channel %s is already exist", channelID)
	case core.ASSETCHANNELID:
		return fmt.Errorf("Channel %s is already exist", channelID)
	default:
		if !util.IsLegalChannelName(channelID) {
			return fmt.Errorf("%s is not a legal channel name", channelID)
		}
		c.managerLock.RLock()
		if util.Contain(c.Managers, channelID) {
			c.managerLock.RUnlock()
			return fmt.Errorf("Channel %s is aleardy exist", channelID)
		}
		c.managerLock.RUnlock()
	}

	err = c.CM.AddTx(tx)
	if err != nil {
		return err
	}

	c.db.WatchChannel(channelID)

	return nil
}

// getChannelManager return the manager of the channel
func (c *Coordinator) getChannelManager(channelID string) (*Manager, error) {
	switch channelID {
	case core.GLOBALCHANNELID:
		return c.GM, nil
	case core.CONFIGCHANNELID:
		return c.CM, nil
	case core.ASSETCHANNELID:
		return c.AM, nil
	default:
		c.managerLock.RLock()
		defer c.managerLock.RUnlock()
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

	if err := c.loadAssetChannel(); err != nil {
		return err
	}

	return nil
}

// loadConfigChannel load the config channel("_config")
func (c *Coordinator) loadConfigChannel() error {
	var err error
	c.CM, err = NewManager(core.CONFIGCHANNELID, c)
	if err != nil {
		return err
	}
	if !c.CM.HasGenesisBlock() {
		log.Info("Creating genesis block of channel _config")
		// create admins, we just config one admin
		admins, err := bc.CreateAdmins()
		if err != nil {
			return err
		}
		gb, err := bc.CreateGenesisBlock(admins)
		if err != nil {
			return err
		}
		// put  admin's pubkey into leveldb
		err = c.CM.db.UpdateSystemAdmin(&bc.Profile{
			Public: true,
			Admins: admins,
		})
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
	c.GM, err = NewManager(core.GLOBALCHANNELID, c)
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
			ChannelID: core.CONFIGCHANNELID,
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

func (c *Coordinator) loadAssetChannel() error {
	var err error
	c.AM, err = NewManager(core.ASSETCHANNELID, c)
	if err != nil {
		return err
	}
	if !c.AM.HasGenesisBlock() {
		log.Infof("Creating genesis block of channel _asset")
		// agb: asset channel genesis block
		agb, err := ac.CreateGenesisBlock([]*ac.Payload{&ac.Payload{}})
		if err != nil {
			return err
		}
		err = c.AM.AddBlock(agb)
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
			c.managerLock.Lock()
			c.Managers[channelID] = manager
			c.managerLock.Unlock()
			log.Infof("loadUserChannel: load channel %s from leveldb", channelID)
		}
	}
	return nil
}

// setConsensus set consensus according to the consensus config
func (c *Coordinator) setConsensus(cfg *config.ConsensusConfig) error {
	// set consensus
	var channels = make(map[string]consensus.Config, 0)
	defaultCfg := consensus.Config{
		Timeout: c.chainCfg.BatchTimeout,
		MaxSize: c.chainCfg.BatchSize,
		Number:  1,
		Resume:  false,
	}
	channels[core.GLOBALCHANNELID] = defaultCfg
	channels[core.CONFIGCHANNELID] = defaultCfg
	channels[core.ASSETCHANNELID] = defaultCfg
	// set consensus of user channels
	c.managerLock.RLock()
	for channelID := range c.Managers {
		channels[channelID] = defaultCfg
	}
	c.managerLock.RUnlock()
	switch cfg.Type {
	case config.SOLO:
		consensus, err := solo.NewConsensus(channels)
		if err != nil {
			return err
		}
		c.Consensus = consensus
	case config.RAFT:
		raftConfig, err := getConfig(cfg.Raft, c.chainCfg)
		if err != nil {
			return err
		}

		consensus, err := raft.NewConsensus(channels, raftConfig)
		if err != nil {
			return nil
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

func getConfig(cRaft config.RaftConfig, cChain *config.BlockChainConfig) (*raft.Config, error) {
	tlsCfg := consensus.TLSConfig{
		Enable:  cRaft.TLS.Enable,
		CA:      cRaft.TLS.CA,
		RawCert: cRaft.TLS.RawCert,
		Key:     cRaft.TLS.Key,
		Pool:    cRaft.TLS.Pool,
		Cert:    cRaft.TLS.Cert,
	}
	return raft.NewConfig(cRaft.Path, cRaft.ID, cRaft.Nodes, cRaft.Join, consensus.Config{
		Timeout: cChain.BatchTimeout,
		MaxSize: cChain.BatchSize,
		Resume:  false,
		Number:  1,
		TLS:     tlsCfg,
	})
}
