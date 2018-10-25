package tendermint

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// Glue will connect consensus and tendermint
type Glue struct {
	lock *sync.Mutex
	// tn is the number(height) of tendermint
	tn int64
	// th is the hash of tendermint
	th     []byte
	txs    [][]byte
	db     *DB
	port   int
	client *Client
}

// NewGlue is the constructor of Glue
func NewGlue(dbDir string, port int) (*Glue, error) {
	g := new(Glue)
	g.lock = new(sync.Mutex)
	db, err := NewDB(dbDir)
	if err != nil {
		return nil, err
	}
	g.db = db
	g.tn = db.GetHeight()
	g.th = db.GetHash()
	g.port = port
	return g, nil
}

// Start run the glue
func (g *Glue) Start() error {
	// Start the listener
	srv, err := server.NewServer(fmt.Sprintf("tcp://0.0.0.0:%d", g.port), "socket", g)
	if err != nil {
		return err
	}
	srv.SetLogger(NewLogger())
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
	return nil
}

// CheckTx always return OK
func (g *Glue) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

// DeliverTx add tx into txs and return OK
func (g *Glue) DeliverTx(tx []byte) types.ResponseDeliverTx {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.txs = append(g.txs, tx)
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

// Commit will generate a block and init the txs
func (g *Glue) Commit() types.ResponseCommit {
	g.lock.Lock()
	defer g.lock.Unlock()

	// todo: how to manage these txs is still consided
	if len(g.txs) != 0 {

	}

	g.db.SetHeight(g.tn)
	g.db.SetHash(g.th)
	return types.ResponseCommit{}
}

// BeginBlock set height and hash
func (g *Glue) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.tn = req.Header.Height
	g.th = req.Header.AppHash
	return types.ResponseBeginBlock{}
}

// EndBlock is not support validator updates now
func (g *Glue) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	g.lock.Lock()
	defer g.lock.Unlock()

	return types.ResponseEndBlock{}
}

// Info is used to avoid load all blocks
func (g *Glue) Info(req types.RequestInfo) types.ResponseInfo {
	g.lock.Lock()
	defer g.lock.Unlock()

	return types.ResponseInfo{LastBlockHeight: g.db.GetHeight(), LastBlockAppHash: g.db.GetHash()}
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
// TODO: Not implementation yet
func (g *Glue) AddTx(channelID string, tx []byte) error {
	NewTx(channelID, tx)
	return nil
}
