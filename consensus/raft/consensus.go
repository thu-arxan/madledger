package raft

import (
	"errors"
	"madledger/consensus"
	"time"

	"google.golang.org/grpc"
)

// Consensus is the implementation of interface
// Consensus keeps connections to raft services
type Consensus struct {
	cfg     *Config
	chain   *BlockChain
	clients []*Client
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

// NewConseneus is the constructor of Consensus
func NewConseneus(cfg *Config) (*Consensus, error) {
	chain, err := NewBlockChain(cfg.cc)
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
	if err := c.chain.Start(); err != nil {
		return err
	}

	c.clients = make([]*Client, 0)
	for _, addr := range c.cfg.ec.peers {
		client, err := NewClient(addr)
		if err != nil {
			return err
		}
		c.clients = append(c.clients, client)
	}

	return errors.New("Not implementation yet")
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	return errors.New("Not implementation yet")
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return errors.New("Not implementation yet")
}

// AddChannel is the implementation of interface
// Note: we can ignore this function now
func (c *Consensus) AddChannel(channelID string, cfg EraftConfig) error {
	return nil
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return nil, errors.New("Not implementation yet")
}
