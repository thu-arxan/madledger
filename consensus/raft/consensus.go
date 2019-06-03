package raft

import (
	"context"
	"errors"
	"madledger/consensus"
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
	clients []*Client
	leader  int32
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
	c.clients = make([]*Client, 0)
	for _, addr := range c.cfg.ec.peers {
		client, err := NewClient(addr)
		if err != nil {
			return err
		}
		c.clients = append(c.clients, client)
	}

	if err := c.chain.Start(); err != nil {
		return err
	}

	return nil
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	return errors.New("Not implementation yet")
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	var err error
	log.Infof("Add tx of channel %s", channelID)
	for i := 0; i < 10; i++ {
		log.Infof("Try to add tx to %d", c.leader)
		err = c.clients[c.getLeader()].addTx(channelID, tx)
		if err == nil {
			return nil
		}
		c.setLeader(util.RandNum(len(c.clients)))
		log.Infof("Retry %d times and leader is %d", i, c.leader)
		time.Sleep(200 * time.Millisecond)
	}

	log.Infof("Add tx failed because %s", err)
	return err
}

// AddChannel is the implementation of interface
// Note: we can ignore this function now
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	return nil
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return nil, errors.New("Not implementation yet")
}

func (c *Consensus) setLeader(leader int) {
	atomic.StoreInt32(&c.leader, int32(leader))
}

func (c *Consensus) getLeader() int {
	leader := atomic.LoadInt32(&c.leader)
	return int(leader)
}
