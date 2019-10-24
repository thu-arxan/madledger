package raft

import (
	"context"
	"madledger/consensus"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"madledger/common/util"
	pb "madledger/consensus/raft/protos"

	"google.golang.org/grpc"
)

// Consensus is the implementation of interface
// Consensus keeps connections to raft services
type Consensus struct {
	cfg     *Config
	chain   *BlockChain
	clients map[uint64]*Client
	ids     []uint64
	leader  uint64
}

// Client is the clients keep connections to blockchain
type Client struct {
	addr string
	conn *grpc.ClientConn
}

// NewClient is the constructor of Client
func NewClient(addr string) (*Client, error) {
	return &Client{
		addr: addr,
	}, nil
}

func (c *Client) newConn() error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) addTx(channelID string, tx []byte) error {
	if c.conn == nil {
		if err := c.newConn(); err != nil {
			return err
		}
	}
	client := pb.NewBlockChainClient(c.conn)
	_, err := client.AddTx(context.Background(), &pb.Tx{
		Data: NewTx(channelID, tx).Bytes(),
	})
	return err
}

// NewConseneus is the constructor of Consensus
func NewConseneus(cfg *Config) (*Consensus, error) {
	chain, err := NewBlockChain(cfg)
	if err != nil {
		return nil, err
	}
	return &Consensus{
		cfg:   cfg,
		chain: chain,
	}, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	c.clients = make(map[uint64]*Client)
	c.ids = make([]uint64, 0)
	for id, addr := range c.cfg.ec.peers {
		client, err := NewClient(addr)
		if err != nil {
			return err
		}
		c.clients[id] = client
		c.ids = append(c.ids, id)
	}
	c.setLeader(c.ids[util.RandNum(len(c.ids))])

	if err := c.chain.Start(); err != nil {
		return err
	}

	return nil
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	c.chain.Stop()
	c.chain.raft.Stop()
	c.chain.rpcServer.Stop()
	return nil
	//return errors.New("Not implementation yet")
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	var err error
	log.Infof("raft.Consensus.AddTx: Add tx of channel %s.", channelID)
	// todo: we should parse the leader address other than random choose a leader
	for i := 0; i < 100; i++ {
		log.Infof("Try to add tx to raft %d, this is %d times' trying.", c.leader, i)
		err = c.clients[c.getLeader()].addTx(channelID, tx)
		if err == nil || strings.Contains(err.Error(), "Transaction is aleardy in the pool") {
			log.Infof("Succeed to add tx to raft %d", c.leader)
			return nil
		}
		log.Info(err)
		// then parse leader id
		if strings.Contains(err.Error(), "Please send to leader") {
			id, err := strconv.ParseUint(strings.Replace(err.Error(), "rpc error: code = Unknown "+
				"desc = Please send to leader ", "", -1), 10, 64)
			if err == nil && id != 0 {
				c.setLeader(id)
				continue
			}
		}else if strings.Contains(err.Error(),"I will stop and can not add tx to chain."){
			return err
		}
		// error except tx exist and the id is not leader
		log.Infoln("error except tx exist and the id is not leader, so set leader randomly.")
		c.setLeader(c.ids[util.RandNum(len(c.ids))])
		time.Sleep(200 * time.Millisecond)
	}

	log.Infof("Add tx failed because %s", err)
	return err
}

// AddChannel is the implementation of interface
// Note: we can ignore this function now
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	log.Infof("Add channel %s", channelID)
	return nil
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	// log.Infof("Get block %d of channel %s", num, channelID)
	return c.chain.getBlock(channelID, num, async)
}

func (c *Consensus) setLeader(leader uint64) {
	atomic.StoreUint64(&c.leader, leader)
}

func (c *Consensus) getLeader() uint64 {
	leader := atomic.LoadUint64(&c.leader)
	return leader
}
