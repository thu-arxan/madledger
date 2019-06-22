package raft

import (
	"context"
	"errors"
	"fmt"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	pb "madledger/consensus/raft/protos"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

// BlockChain will create blockchain in raft
type BlockChain struct {
	lock sync.Mutex

	raft          *Raft
	config        consensus.Config
	txs           chan bool
	pool          *txPool
	hub           *event.Hub
	num           uint64
	//hybridBlocks  map[uint64]*HybridBlock
	hybridBlockCh chan *HybridBlock

	rpcServer *grpc.Server

	quit chan bool
	done chan bool
}

// NewBlockChain is the constructor of blockchain
func NewBlockChain(cfg *Config) (*BlockChain, error) {
	raft, err := NewRaft(cfg.ec)
	if err != nil {
		return nil, err
	}

	return &BlockChain{
		config:        cfg.cc,
		txs:           make(chan bool, cfg.cc.MaxSize),
		quit:          make(chan bool, 1),
		done:          make(chan bool, 1),
		hybridBlockCh: raft.BlockCh(),
		pool:          newTxPool(),
		raft:          raft,
		hub:           event.NewHub(),
		//hybridBlocks:  make(map[uint64]*HybridBlock),
	}, nil
}

// Start start the blockchain service
func (chain *BlockChain) Start() error {
	if err := chain.raft.Start(); err != nil {
		return err
	}
	atomic.StoreUint64(&(chain.num), chain.raft.app.db.GetChainNum())

	if err := chain.start(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", chain.raft.cfg.getLocalChainAddress())
	if err != nil {
		return fmt.Errorf("Failed to start the server(%s)", err)
	}
	log.Infof("Listen %s", chain.raft.cfg.getLocalChainAddress())
	var opts []grpc.ServerOption
	chain.rpcServer = grpc.NewServer(opts...)
	pb.RegisterBlockChainServer(chain.rpcServer, chain)
	go func() {
		err = chain.rpcServer.Serve(lis)
		if err != nil {
			log.Error("Start server failed: ", err)
			return
		}
	}()

	time.Sleep(300 * time.Millisecond)
	return nil
}

// start chain service
func (chain *BlockChain) start() error {
	log.Infof("Raft blockchain start")

	go func() {
		ticker := time.NewTicker(time.Duration(chain.config.Timeout) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
			case <-chain.txs:
				if chain.pool.getPoolSize() >= chain.config.MaxSize {
					chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
				}
			case block := <-chain.hybridBlockCh:
				if block.GetNumber() == chain.num+1 {
					log.Infof("Blockchain.start: call addBlock and the block number is %d", block.GetNumber())
					chain.addBlock(block)
				} else if block.GetNumber() <= chain.num {
					chain.raft.FetchBlockDone(block.GetNumber())
				} else {
					chain.raft.NotifyLater(block)
				}
			case <-chain.quit:
				chain.done <- true
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

// AddTx will try to add a tx
func (chain *BlockChain) AddTx(ctx context.Context, in *pb.Tx) (*pb.None, error) {
	err := chain.addTx(in.Data)
	return &pb.None{}, err
}

func (chain *BlockChain) addTx(tx []byte) error {
	if !chain.raft.IsLeader() {
		return fmt.Errorf("Please send to leader %d", chain.raft.GetLeader())
	}
	log.Infof("[%d] add tx", chain.raft.cfg.id)
	err := chain.pool.addTx(tx)

	if err != nil {
		return err
	}

	go func() {
		chain.txs <- true
	}()

	log.Infof("[%d]Hub watch tx %s", chain.raft.cfg.id, util.Hex(crypto.Hash(tx)))
	result := chain.hub.Watch(util.Hex(crypto.Hash(tx)), nil)
	return result.Err
}

// Stop will block the work of channel
func (chain *BlockChain) Stop() {
	chain.quit <- true
	<-chain.done
}

func (chain *BlockChain) createBlock(txs [][]byte) error {
	if len(txs) == 0 {
		return nil
	}
	block := &HybridBlock{
		Num: chain.num + 1,
		Txs: txs,
	}
	// then call eraft
	log.Infof("[%d]Try to add block %d", chain.raft.cfg.id, block.Num)
	if err := chain.raft.AddBlock(block); err != nil {
		// todo: if we failed to create block we should release all txs
		log.Infof("[%d]Failed to add block %d because %s", chain.raft.cfg.id, block.Num, err)
		return err
	}

	log.Infof("[%d]Succeed to add block %d", chain.raft.cfg.id, block.Num)
	log.Infof("Blockchain.createBlock: call addBlock adn the block number is %d", block.GetNumber())
	chain.addBlock(block)
	return nil
}

func (chain *BlockChain) addBlock(block *HybridBlock) {
	//chain.hybridBlocks[block.GetNumber()] = block
	chain.num = block.GetNumber()
	chain.raft.app.db.SetChainNum(block.GetNumber())
	var txs = make(map[string][][]byte)
	for _, tx := range block.Txs {
		hash := util.Hex(crypto.Hash(tx))
		log.Infof("[%d]Hub done tx %s", chain.raft.cfg.id, hash)
		chain.hub.Done(hash, nil)
		t, _ := UnmarshalTx(tx)
		if !util.Contain(txs, t.ChannelID) {
			txs[t.ChannelID] = make([][]byte, 0)
		}
		txs[t.ChannelID] = append(txs[t.ChannelID], t.Data)
	}
	chain.hub.Done(string(block.GetNumber()), nil)
	for channel := range txs {
		num := chain.raft.app.db.GetPrevBlockNum(channel) + 1
		block := &Block{
			ChannelID: channel,
			Num:       num,
			Txs:       txs[channel],
		}
		chain.raft.app.db.AddBlock(block)
		chain.raft.app.db.SetPrevBlockNum(channel, num)
	}
}

func (chain *BlockChain) getBlock(channelID string, num uint64, async bool) (*Block, error) {
	block := chain.raft.app.db.GetBlock(channelID, num, async)
	if block != nil {
		return block, nil
	}
	return nil, errors.New("Not exist")
}
