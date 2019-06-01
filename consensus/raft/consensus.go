package raft

import (
	"context"
	"errors"
	"madledger/consensus"
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
	leader  int
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
	conn, err := grpc.Dial(c.addr, nil, grpc.WithTimeout(2000*time.Millisecond))
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
	log.Info("Consensus Start function")
	if err := c.chain.Start(); err != nil {
		return err
	}

	log.Info("Create clients...")
	c.clients = make([]*Client, 0)
	for _, addr := range c.cfg.ec.peers {
		client, err := NewClient(addr)
		if err != nil {
			return err
		}
		c.clients = append(c.clients, client)
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
	for i := 0; i < 10; i++ {
		err = c.clients[c.leader].addTx(channelID, tx)
		if err == nil {
			return nil
		}
		c.leader = util.RandNum(len(c.clients))
		time.Sleep(100 * time.Millisecond)
	}

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
