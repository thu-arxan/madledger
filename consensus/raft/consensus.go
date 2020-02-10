package raft

import (
	"context"
	"crypto/tls"
	"madledger/consensus"
	"madledger/core"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/credentials"

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
	TLS  consensus.TLSConfig
}

// NewClient is the constructor of Client
func NewClient(addr string, tlsConfig consensus.TLSConfig) (*Client, error) {
	return &Client{
		addr: addr,
		TLS:  tlsConfig,
	}, nil
}

func (c *Client) newConn() error {
	var opts []grpc.DialOption
	var conn *grpc.ClientConn
	var err error
	if c.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*(c.TLS.Cert)},
			RootCAs:      c.TLS.Pool,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(2000*time.Millisecond))
	conn, err = grpc.Dial(c.addr, opts...)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) addTx(channelID string, tx []byte, caller uint64) error {
	if c.conn == nil {
		if err := c.newConn(); err != nil {
			return err
		}
	}
	client := pb.NewBlockChainClient(c.conn)
	_, err := client.AddTx(context.Background(), &pb.RaftTX{
		Tx:     NewTx(channelID, tx).Bytes(),
		Caller: caller,
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
	for id, addr := range c.cfg.ec.GetPeers() {
		client, err := NewClient(addr, c.cfg.cc.TLS)
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
	// c.chain.raft.Stop()
	// c.chain.rpcServer.Stop()
	return nil
	//return errors.New("Not implementation yet")
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(tx *core.Tx) error {
	var err error

	channelID := tx.Data.ChannelID
	bytes, _ := tx.Bytes()
	log.Infof("Consensus.AddTx: Add tx of channel %s.", channelID)
	// todo: we should parse the leader address other than random choose a leader
	for i := 0; i < 100; i++ {
		log.Infof("Try to add tx to raft %d, this is %d times' trying.", c.leader, i)
		err = c.clients[c.getLeader()].addTx(channelID, bytes, c.cfg.id)
		if err == nil || strings.Contains(err.Error(), "Transaction is aleardy in the pool") {
			log.Infof("Succeed to add tx to raft %d, I'm raft %d", c.leader, c.cfg.id)
			return nil
		}

		log.Debugf("add tx meets error:%v", err)
		// then parse leader id
		if strings.Contains(err.Error(), "Please send to leader") {
			id, err := strconv.ParseUint(strings.Replace(err.Error(), "rpc error: code = Unknown "+
				"desc = Please send to leader ", "", -1), 10, 64)
			if err == nil && id != 0 {
				c.setLeader(id)
				continue
			}
		}
		if strings.Contains(err.Error(), "I've been removed from cluster") {
			return err
		}
		// error except tx exist and the id is not leader
		log.Info("Error unknown, set leader randomly.")
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
