package tendermint

import (
	"encoding/json"
	"fmt"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"sync"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// Glue will connect consensus and tendermint
type Glue struct {
	lock *sync.Mutex
	// tn is the number(height) of tendermint
	tn int64
	// th is the hash of tendermint
	th  []byte
	txs [][]byte

	hub     *event.Hub
	blocks  map[string][]*Block
	chans   map[string]*chan consensus.Block
	dbDir   string
	db      *DB
	port    int
	rpcPort int
	client  *Client

	srv cmn.Service
}

// NewGlue is the constructor of Glue
func NewGlue(dbDir string, port *Port) (*Glue, error) {
	g := new(Glue)
	g.lock = new(sync.Mutex)

	g.dbDir = dbDir
	g.blocks = make(map[string][]*Block)
	g.chans = make(map[string]*chan consensus.Block)
	g.port = port.App
	g.rpcPort = port.RPC
	return g, nil
}

// Start run the glue
func (g *Glue) Start() error {
	db, err := NewDB(g.dbDir)
	if err != nil {
		return fmt.Errorf("Failed to load db at %s because %s", g.dbDir, err.Error())
	}
	g.db = db
	g.tn = db.GetHeight()
	g.th = db.GetHash()
	g.hub = event.NewHub()

	g.client, err = NewClient(g.rpcPort)
	if err != nil {
		return err
	}
	log.Info("Start glue...")

	return nil
}

// Stop stop the glue service
// TODO: This way may be too violent
func (g *Glue) Stop() {
	g.db.Close()
	log.Info("Succeed to stop the glue")
}

// CheckTx always return OK
func (g *Glue) CheckTx(tx []byte) types.ResponseCheckTx {
	//t, _ := BytesToTx(tx)

	//log.Infof("[%d]Check Tx %s", g.port, string(t.Data))
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

// DeliverTx add tx into txs and return OK
func (g *Glue) DeliverTx(tx []byte) types.ResponseDeliverTx {
	//t, _ := BytesToTx(tx)

	//log.Infof("[%d]Deliever Tx %s", g.port, string(t.Data))
	g.lock.Lock()
	defer g.lock.Unlock()

	g.txs = append(g.txs, tx)
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

// Commit will generate a block and init the txs
func (g *Glue) Commit() types.ResponseCommit {
	log.Infof("[%d]Commit at block %d", g.port, g.tn)
	g.lock.Lock()
	defer g.lock.Unlock()

	if len(g.txs) != 0 {
		var txs = make(map[string][][]byte)
		for i := range g.txs {
			tx, err := BytesToTx(g.txs[i])
			if err == nil {
				if !util.Contain(txs, tx.ChannelID) {
					txs[tx.ChannelID] = make([][]byte, 0)
				}
				txs[tx.ChannelID] = append(txs[tx.ChannelID], tx.Data)
			}
		}
		for channelID := range txs {
			log.Infof("This is range of channel %s", channelID)
			var num uint64
			if !util.Contain(g.blocks, channelID) {
				g.blocks[channelID] = make([]*Block, 0)
				num = g.db.GetChannelBlockNumber(channelID) + 1
			} else {
				if len(g.blocks[channelID]) != 0 {
					num = g.blocks[channelID][len(g.blocks[channelID])-1].GetNumber() + 1
				}
			}
			block := &Block{
				ChannelID: channelID,
				Num:       num,
				Txs:       txs[channelID],
			}
			g.blocks[channelID] = append(g.blocks[channelID], block)
			g.db.AddBlock(block)
			log.Infof("Done block %s", fmt.Sprintf("%s:%d", channelID, num))
			g.hub.Done(fmt.Sprintf("%s:%d", channelID, num), nil)
			// todo: if we haven't set sync channel, here will lost the block
			go func(channelID string) {
				log.Infof("[%d]Send block of channel %s:%d", g.port, channelID, block.GetNumber())
				if util.Contain(g.chans, channelID) {
					(*g.chans[channelID]) <- block
				}
			}(channelID)
		}
		g.txs = make([][]byte, 0)
	}

	g.db.SetHeight(g.tn)
	g.db.SetHash(g.th)

	log.Infof("[%d]Commit at block %d done", g.port, g.tn)
	return types.ResponseCommit{}
}

// BeginBlock set height and hash
func (g *Glue) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	log.Infof("[%d]Begin block %d", g.port, req.Header.Height)
	g.lock.Lock()
	defer g.lock.Unlock()

	g.tn = req.Header.Height
	g.th = req.Header.AppHash
	return types.ResponseBeginBlock{}
}

// EndBlock is not support validator updates now
func (g *Glue) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	log.Infof("[%d]End block %d", g.port, g.tn)
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Infof("[%d]End block %d done", g.port, g.tn)
	return types.ResponseEndBlock{}
}

// Info is used to avoid load all blocks
func (g *Glue) Info(req types.RequestInfo) types.ResponseInfo {
	log.Infof("[%d]Info", g.port)
	g.lock.Lock()
	defer g.lock.Unlock()

	return types.ResponseInfo{
		LastBlockHeight:  g.db.GetHeight(),
		LastBlockAppHash: g.db.GetHash(),
	}
}

// InitChain just send the init chain message
func (g *Glue) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	g.lock.Lock()
	defer g.lock.Unlock()

	return types.ResponseInitChain{}
}

// SetOption is not useful now
func (g *Glue) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	g.lock.Lock()
	defer g.lock.Unlock()

	return types.ResponseSetOption{}
}

// Query is not working now
func (g *Glue) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	return types.ResponseQuery{}
}

// SetSyncChan set the sync chan of channelID
func (g *Glue) SetSyncChan(channelID string, ch *chan consensus.Block) {
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Infof("Set sync chan of channel %s", channelID)
	g.chans[channelID] = ch
}

// GetBlock return the block of channelID and with the num
func (g *Glue) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	g.lock.Lock()
	for i := range g.blocks[channelID] {
		if g.blocks[channelID][i].GetNumber() == num {
			defer g.lock.Unlock()
			log.Infof("consensus/tendermint/app: get block %d from g.blocks[%s]", num, channelID)
			return g.blocks[channelID][i], nil
		}
	}
	g.lock.Unlock()
	// But block is not in blocks does not mean it is not exist
	block := g.db.GetBlock(channelID, num)
	if block != nil {
		log.Infof("consensus/tendermint/app: get block %d from g.db and key is %s", num, channelID)
		return block, nil
	}

	if async {
		log.Infof("Watch block %s", fmt.Sprintf("%s:%d", channelID, num))
		g.hub.Watch(fmt.Sprintf("%s:%d", channelID, num), nil)
		g.lock.Lock()
		defer g.lock.Unlock()
		for i := range g.blocks[channelID] {
			if g.blocks[channelID][i].GetNumber() == num {
				log.Infof("consensus/tendermint/app: get block %d from g.blocks[%s] asynchronously", num, channelID)
				return g.blocks[channelID][i], nil
			}
		}
	}

	return nil, fmt.Errorf("Block %s:%d is not exist", channelID, num)
}

// Tx is the union of ChannelID and Data
type Tx struct {
	ChannelID string
	Data      []byte
}

// NewTx is the constructor of Tx
func NewTx(channelID string, data []byte) *Tx {
	return &Tx{
		ChannelID: channelID,
		Data:      data,
	}
}

// BytesToTx convert bytes to Tx
func BytesToTx(bs []byte) (*Tx, error) {
	var t Tx
	err := json.Unmarshal(bs, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Bytes return bytes of Tx
func (t *Tx) Bytes() []byte {
	bs, _ := json.Marshal(t)
	return bs
}

// AddTx add a tx
func (g *Glue) AddTx(channelID string, tx []byte) error {
	//log.Infof("[%d]Channel %s add tx %s", g.port, channelID, string(tx))
	return g.client.AddTx(NewTx(channelID, tx).Bytes())
}
