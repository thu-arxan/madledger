package raft

import (
	"context"
	"crypto/tls"
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
	pb "madledger/consensus/raft/protos"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// BlockChain will create blockchain in raft
// BlockChain handles grpc requests: AddTx, etc.
// BlockChain manages channels which will packet txs separately
type BlockChain struct {
	sync.RWMutex
	cfg       *Config
	channels  map[string]*channel // channel id => channel
	raft      *eraft.Raft
	rpcServer *grpc.Server
	stop      chan *chan bool
}

// NewBlockChain is the constructor of blockchain
func NewBlockChain(cfg *Config) (*BlockChain, error) {
	raft, err := eraft.NewRaft(cfg.ec)
	if err != nil {
		return nil, err
	}
	// todo: more config
	// todo: load channels, setting
	return &BlockChain{
		cfg:      cfg,
		raft:     raft,
		channels: make(map[string]*channel),
		stop:     make(chan *chan bool),
	}, nil
}

// Start start the blockchain service
func (chain *BlockChain) Start() error {
	// start raft
	if err := chain.raft.Start(); err != nil {
		return err
	}

	// start channels
	chain.RLock()
	for _, channel := range chain.channels {
		channel.start()
	}
	chain.RUnlock()

	// start grpc service
	addr := chain.cfg.peers[chain.cfg.id]
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen failed: %v", err)
	}

	log.Infof("raft.chain[%d] listen on: %s", chain.cfg.id, addr)

	var opts []grpc.ServerOption
	if chain.cfg.cc.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{*(chain.cfg.cc.TLS.Cert)},
			ClientCAs:    chain.cfg.cc.TLS.Pool,
		})
		opts = append(opts, grpc.Creds(creds))
	}
	chain.rpcServer = grpc.NewServer(opts...)
	pb.RegisterBlockChainServer(chain.rpcServer, chain)
	go func() {
		err = chain.rpcServer.Serve(ln)
		if err != nil {
			log.Errorf("Start server failed: %v", err)
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Stop will block the work of channel
func (chain *BlockChain) Stop() {
	chain.Lock()
	defer chain.Unlock()
	// todo: change Stop to GracefulStop?
	if chain.rpcServer != nil {
		chain.rpcServer.Stop()
	}

	for _, channel := range chain.channels {
		channel.Stop()
	}

	chain.raft.Stop()
}

// AddTx will try to add a tx
func (chain *BlockChain) AddTx(ctx context.Context, in *pb.RaftTX) (*pb.None, error) {
	if chain.raft.Removed(in.Caller) {
		log.Infof("[%d]I've been removed from cluster.", in.Caller)
		return &pb.None{}, fmt.Errorf("[%d]I've been removed from cluster", in.Caller)
	}
	if !chain.raft.IsLeader() {
		// get leader and return
		return &pb.None{}, fmt.Errorf("%s %d", NotLeaderMsg, chain.raft.GetLeader())
	}
	channel, err := chain.getChannel(in.Channel)
	if err != nil {
		return &pb.None{}, err
	}
	err = channel.addTx(in.Tx)
	return &pb.None{}, err
}

func (chain *BlockChain) addChannels(channelID string, cfg consensus.Config) error {
	chain.Lock()
	defer chain.Unlock()

	if util.Contain(chain.channels, channelID) {
		return fmt.Errorf("channel %s exits", channelID)
	}

	channel := newChannel(chain.cfg.id, channelID, cfg, chain.raft)
	chain.channels[channelID] = channel
	log.Infof("add channel %s succeed", channelID)
	return nil
}

func (chain *BlockChain) startChannel(channelID string) error {
	channel, err := chain.getChannel(channelID)
	if err != nil {
		return err
	}
	return channel.start()
}

func (chain *BlockChain) getChannel(channelID string) (*channel, error) {
	chain.RLock()
	defer chain.RUnlock()

	if util.Contain(chain.channels, channelID) {
		return chain.channels[channelID], nil
	}

	return nil, fmt.Errorf("channel %s not exist", channelID)
}

func (chain *BlockChain) getBlock(channelID string, num uint64, async bool) (*eraft.Block, error) {
	channel, err := chain.getChannel(channelID)
	if err != nil {
		return nil, err
	}

	return channel.getBlock(num, async)
}
