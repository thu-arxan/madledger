package raft

import (
	"context"
	"errors"
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
	pb "madledger/consensus/raft/protos"
	"sync"
)

// import (
// 	"context"
// 	"crypto/tls"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"madledger/common"
// 	"madledger/common/crypto"
// 	"madledger/common/event"
// 	"madledger/common/util"
// 	"madledger/consensus"
// 	"madledger/consensus/raft/eraft"
// 	"madledger/core"

// 	"net"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"go.etcd.io/etcd/pkg/types"
// 	"go.etcd.io/etcd/raft/raftpb"
// 	"google.golang.org/grpc/credentials"

// 	"google.golang.org/grpc"
// )

// BlockChain will create blockchain in raft
type BlockChain struct {
	sync.RWMutex

	config   consensus.Config
	channels map[string]*channel // channel id => channel
	// raft   *eraft.Raft

	// txs    chan bool
	// pool   *txPool
	// hub    *event.Hub
	// num    uint64
	// //hybridBlocks  map[uint64]*HybridBlock
	// hybridBlockCh chan *eraft.HybridBlock
	// // record confChange
	// cfgChanges []raftpb.ConfChange

	// rpcServer *grpc.Server

	// quit chan bool
	// done chan bool
}

// NewBlockChain is the constructor of blockchain
func NewBlockChain(cfg *Config) (*BlockChain, error) {
	// todo: more config
	// todo: load channels, setting
	return &BlockChain{
		channels: make(map[string]*channel),
	}, nil
}

// // NewBlockChain is the constructor of blockchain
// func NewBlockChain(cfg *Config) (*BlockChain, error) {
// 	raft, err := eraft.NewRaft(cfg.ec)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &BlockChain{
// 		config:        cfg.cc,
// 		txs:           make(chan bool, cfg.cc.MaxSize),
// 		quit:          make(chan bool, 1),
// 		done:          make(chan bool, 1),
// 		hybridBlockCh: raft.BlockCh(),
// 		pool:          newTxPool(),
// 		raft:          raft,
// 		hub:           event.NewHub(),
// 		//hybridBlocks:  make(map[uint64]*HybridBlock),
// 		cfgChanges: make([]raftpb.ConfChange, 0),
// 	}, nil
// }

// Start start the blockchain service
func (chain *BlockChain) Start() error {
	// todo: start channels
	// start grpc server, etc
	chain.Lock()
	defer chain.Unlock()
	for _, channel := range chain.channels {
		channel.start()
	}
	return nil
}

// Stop will block the work of channel
func (chain *BlockChain) Stop() {
	// todo: close grpc server, clean setting, etc
}

// // Start start the blockchain service
// func (chain *BlockChain) Start() error {
// 	if err := chain.raft.Start(); err != nil {
// 		return err
// 	}
// 	atomic.StoreUint64(&(chain.num), chain.raft.app.db.GetChainNum())

// 	if err := chain.start(); err != nil {
// 		return err
// 	}

// 	lis, err := net.Listen("tcp", chain.raft.cfg.getLocalChainAddress())
// 	if err != nil {
// 		return fmt.Errorf("Failed to start the server, because %s)", err)
// 	}
// 	log.Infof("Listen %s", chain.raft.cfg.getLocalChainAddress())
// 	var opts []grpc.ServerOption
// 	if chain.config.TLS.Enable {
// 		creds := credentials.NewTLS(&tls.Config{
// 			ClientAuth:   tls.RequireAndVerifyClientCert,
// 			Certificates: []tls.Certificate{*(chain.config.TLS.Cert)},
// 			ClientCAs:    chain.config.TLS.Pool,
// 		})
// 		opts = append(opts, grpc.Creds(creds))
// 	}
// 	chain.rpcServer = grpc.NewServer(opts...)
// 	pb.RegisterBlockChainServer(chain.rpcServer, chain)
// 	go func() {
// 		err = chain.rpcServer.Serve(lis)
// 		if err != nil {
// 			log.Error("Start server failed: ", err)
// 			return
// 		}
// 	}()

// 	time.Sleep(300 * time.Millisecond)
// 	return nil
// }

// // start chain service
// func (chain *BlockChain) start() error {
// 	log.Infof("Raft blockchain start")

// 	go func() {
// 		ticker := time.NewTicker(time.Duration(chain.config.Timeout) * time.Millisecond)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-ticker.C:
// 				chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
// 			case <-chain.txs:
// 				if chain.pool.getPoolSize() >= chain.config.MaxSize {
// 					chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
// 				}
// 			case block := <-chain.hybridBlockCh:
// 				if block.GetNumber() == chain.num+1 {
// 					log.Infof("Blockchain.start: call addBlock and the block number is %d", block.GetNumber())
// 					chain.addBlock(block)
// 				} else if block.GetNumber() <= chain.num {
// 					chain.raft.FetchBlockDone(block.GetNumber())
// 				} else {
// 					chain.raft.NotifyLater(block)
// 				}
// 			case <-chain.quit:
// 				chain.done <- true
// 				return
// 			}
// 		}
// 	}()

// 	time.Sleep(100 * time.Millisecond)
// 	return nil
// }

// AddTx will try to add a tx
func (chain *BlockChain) AddTx(ctx context.Context, in *pb.RaftTX) (*pb.None, error) {
	// todo: check leader, add to channedl
	if !chain.isLeader() {
		// get leader and return
		return &pb.None{}, errors.New("not leader")
	}
	channel, err := chain.getChannel(in.Channel)
	if err != nil {
		return &pb.None{}, err
	}
	err = channel.addTx(in.Tx)
	return &pb.None{}, err
}

func (chain *BlockChain) isLeader() bool {
	// todo: check leader
	return false
}

func (chain *BlockChain) getChannel(name string) (*channel, error) {
	// todo
	return nil, errors.New("not implement yet")
}

// func (chain *BlockChain) addTx(tx []byte, caller uint64) error {
// 	// check if the node has been removed
// 	if chain.raft.eraft.removed[types.ID(caller)] {
// 		log.Infof("[%d]I've been removed from cluster.", caller)
// 		return fmt.Errorf("[%d]I've been removed from cluster", caller)
// 	}
// 	if !chain.raft.IsLeader() {
// 		return fmt.Errorf("Please send to leader %d", chain.raft.GetLeader())
// 	}

// 	log.Infof("[%d]add tx", chain.raft.cfg.id)
// 	err := chain.pool.addTx(tx)
// 	if err != nil {
// 		return err
// 	}

// 	go func() {
// 		chain.txs <- true
// 	}()

// 	log.Infof("[%d]Hub watch tx %s", chain.raft.cfg.id, util.Hex(crypto.Hash(tx)))
// 	result := chain.hub.Watch(util.Hex(crypto.Hash(tx)), nil)
// 	return result.Err
// }

// // Stop will block the work of channel
// func (chain *BlockChain) Stop() {
// 	chain.quit <- true
// 	<-chain.done
// }

// func (chain *BlockChain) createBlock(txs [][]byte) error {
// 	if len(txs) == 0 {
// 		return nil
// 	}
// 	block := &eraft.HybridBlock{
// 		Num: chain.num + 1,
// 		Txs: txs,
// 	}
// 	// then call eraft
// 	log.Infof("[%d]Try to add block %d", chain.raft.cfg.id, block.Num)
// 	if err := chain.raft.AddBlock(block); err != nil {
// 		// todo: if we failed to create block we should release all txs
// 		log.Infof("[%d]Failed to add block %d because %s", chain.raft.cfg.id, block.Num, err)
// 		return err
// 	}

// 	log.Infof("[%d]Succeed to add block %d", chain.raft.cfg.id, block.Num)
// 	log.Infof("Blockchain.createBlock: call addBlock and the block number is %d", block.GetNumber())
// 	return chain.addBlock(block)
// }

// func (chain *BlockChain) addBlock(block *eraft.HybridBlock) error {
// 	//chain.hybridBlocks[block.GetNumber()] = block
// 	chain.num = block.GetNumber()
// 	chain.raft.app.db.SetChainNum(block.GetNumber())
// 	var txs = make(map[string][][]byte)
// 	var err error
// 	for _, tx := range block.Txs {
// 		hash := util.Hex(crypto.Hash(tx))
// 		log.Infof("[%d]Hub done tx %s", chain.raft.cfg.id, hash)
// 		chain.hub.Done(hash, nil)
// 		t, _ := UnmarshalTx(tx)
// 		if !util.Contain(txs, t.ChannelID) {
// 			txs[t.ChannelID] = make([][]byte, 0)
// 		}
// 		txs[t.ChannelID] = append(txs[t.ChannelID], t.Data)
// 		// every node need get confChange to avoid leader change
// 		cfgChange, err := chain.getConfChange(t)
// 		if err != nil {
// 			if !strings.Contains(err.Error(), "It's not raft config change tx") {
// 				log.Infof("addBlock meets error:%v", err.Error())
// 				return err
// 			}
// 		} else {
// 			log.Printf("Append cfgChange into cfgChanges, id: %d, type: %v, nodeId: %d, context: %s,"+
// 				" I'm raft %d", cfgChange.ID, cfgChange.Type, cfgChange.NodeID, string(cfgChange.Context), chain.raft.cfg.id)
// 			chain.cfgChanges = append(chain.cfgChanges, cfgChange)
// 		}

// 	}
// 	chain.hub.Done(string(block.GetNumber()), nil)
// 	for channel := range txs {
// 		num := chain.raft.app.db.GetPrevBlockNum(channel) + 1
// 		block := &Block{
// 			ChannelID: channel,
// 			Num:       num,
// 			Txs:       txs[channel],
// 		}
// 		chain.raft.app.db.AddBlock(block)
// 		chain.raft.app.db.SetPrevBlockNum(channel, num)
// 		// propose confChange when the channel is global
// 		// so the removed node can log global channel's block data
// 		if channel == "_global" && len(chain.cfgChanges) > 0 {
// 			if chain.raft.IsLeader() {
// 				for _, change := range chain.cfgChanges {
// 					err = chain.raft.eraft.proposeConfChange(change)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 			}
// 			// every node need clear the cfgChanges if the channel is global
// 			chain.cfgChanges = make([]raftpb.ConfChange, 0)
// 		}
// 	}

// 	return nil
// }

// func (chain *BlockChain) getConfChange(tx *Tx) (raftpb.ConfChange, error) {
// 	var coreTx core.Tx
// 	var cfgChange raftpb.ConfChange
// 	err := json.Unmarshal(tx.Data, &coreTx)
// 	// Note: The reason return nil because Tx may be just random bytes,
// 	// and this is a bad implementation so we should change the way to do this
// 	// TODO: Reimplement it
// 	if err != nil {
// 		return cfgChange, nil
// 	}
// 	// get tx type according to recipient
// 	//log.Infof("Recipient: %s", common.BytesToAddress(coreTx.Data.Recipient).String())
// 	txType, err := core.GetTxType(common.BytesToAddress(coreTx.Data.Recipient).String())
// 	if err == nil && txType == core.NODE {
// 		err = json.Unmarshal(coreTx.Data.Payload, &cfgChange)
// 		if err != nil {
// 			return cfgChange, err
// 		}
// 	} else {
// 		// todo: fucking use illegal words as error
// 		return cfgChange, errors.New("It's not raft config change tx")
// 	}
// 	return cfgChange, err
// }

func (chain *BlockChain) getBlock(channelID string, num uint64, async bool) (*eraft.Block, error) {
	// todo: implement it
	// block := chain.raft.app.db.GetBlock(channelID, num, async)
	// if block != nil {
	// 	return block, nil
	// }
	return nil, errors.New("Not exist")
}
