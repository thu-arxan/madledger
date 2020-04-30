package channel

import (
	"encoding/json"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/consensus"
	"madledger/core"
	"madledger/orderer/config"
	"madledger/orderer/db"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// here we mock a consensus
type mockConsensus struct {
}

func (mc *mockConsensus) Start() error {
	return nil
}

func (mc *mockConsensus) Stop() error {
	return nil
}

func (mc *mockConsensus) AddTx(tx *core.Tx) error {
	return nil
}

func (mc *mockConsensus) AddChannel(channelID string, cfg consensus.Config) error {
	return nil
}

func (mc *mockConsensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	if num != 1 && channelID != core.GLOBALCHANNELID { // global has two blocks
		time.Sleep(10 * time.Second)
	}
	if num > 2 { //must be global channel
		time.Sleep(10 * time.Second)
	}
	switch channelID {
	case core.GLOBALCHANNELID:
		// construct some special global block
		var txs = make([]*core.Tx, 0)
		if num == 1 {
			payload, _ := json.Marshal(core.GlobalTxPayload{
				ChannelID: core.CONFIGCHANNELID,
				Number:    1,
			})
			txs = append(txs, &core.Tx{
				ID: "global1",
				Data: core.TxData{
					ChannelID: core.GLOBALCHANNELID,
					Payload:   payload,
				},
			})
		} else {
			payload, _ := json.Marshal(core.GlobalTxPayload{
				ChannelID: "test",
				Number:    1,
			})
			txs = append(txs, &core.Tx{
				ID: "global2",
				Data: core.TxData{
					ChannelID: core.GLOBALCHANNELID,
					Payload:   payload,
				},
			})
		}
		return &mockBlock{
			num: num,
			txs: txs,
		}, nil

	case core.CONFIGCHANNELID:
		var txs = make([]*core.Tx, 0)
		payload, _ := json.Marshal(cc.Payload{
			ChannelID: "test",
		})
		txs = append(txs, &core.Tx{
			ID: "config1",
			Data: core.TxData{
				ChannelID: core.CONFIGCHANNELID,
				Payload:   payload,
			},
		})
		// then we want insert 10 channels to make sure that the channel update is not to fast
		for i := 0; i < 10; i++ {
			payload, _ := json.Marshal(cc.Payload{
				ChannelID: fmt.Sprintf("test%d", i),
			})
			txs = append(txs, &core.Tx{
				ID: "config1",
				Data: core.TxData{
					ChannelID: core.CONFIGCHANNELID,
					Payload:   payload,
				},
			})
		}
		return &mockBlock{
			num: 1,
			txs: txs,
		}, nil
	case core.ASSETCHANNELID:
	default: //user channel
		var txs = make([]*core.Tx, 0)
		txs = append(txs, &core.Tx{
			ID: "asset1",
			Data: core.TxData{
				ChannelID: channelID,
			},
		})
		return &mockBlock{
			num: 1,
			txs: txs,
		}, nil
	}
	return newMockBlock(num), nil
}

func (mc *mockConsensus) Info() string {
	return "mock"
}

type mockBlock struct {
	num uint64
	txs []*core.Tx
}

func newMockBlock(num uint64) *mockBlock {
	return &mockBlock{
		num: num,
	}
}

func (block *mockBlock) GetTxs() []*core.Tx {
	return block.txs
}

func (block *mockBlock) GetNumber() uint64 {
	return block.num
}

// mockDB is a wrapper of real db, but i want to do something special
type mockDB struct {
	db db.DB
}

func (db *mockDB) ListChannel() []string {
	return db.db.ListChannel()
}

// HasChannel return if channel exist
func (db *mockDB) HasChannel(channelID string) bool {
	return db.db.HasChannel(channelID)
}

// GetChannelProfile return profile of channel if exist, else return an error
func (db *mockDB) GetChannelProfile(channelID string) (*cc.Profile, error) {
	return db.db.GetChannelProfile(channelID)
}

// HasTx return if a tx exist
func (db *mockDB) HasTx(tx *core.Tx) bool {
	return db.db.HasTx(tx)
}

// IsMember return if member belong to channel
func (db *mockDB) IsMember(channelID string, member *core.Member) bool {
	return db.db.IsMember(channelID, member)
}

// IsAdmin return if member if the admin of channel
func (db *mockDB) IsAdmin(channelID string, member *core.Member) bool {
	return db.db.IsAdmin(channelID, member)
}

// GetConsensusBlock return the last consensus block num that db knows
func (db *mockDB) GetConsensusBlock(channelID string) (num uint64) {
	return db.db.GetConsensusBlock(channelID)
}

// WatchChannel provide a way to spy channel change. Now it mainly used to
// spy channel create operation.
func (db *mockDB) WatchChannel(channelID string) {
	db.db.WatchChannel(channelID)
}

// Close close db
func (db *mockDB) Close() error {
	return db.db.Close()
}

// if couldBeEmpty set to true and error is ErrNotFound
// return no error
func (db *mockDB) Get(key []byte, couldBeEmpty bool) ([]byte, error) {
	return db.db.Get(key, couldBeEmpty)
}

// GetAssetAdminPKBytes return nil is not exist
func (db *mockDB) GetAssetAdminPKBytes() []byte {
	return db.db.GetAssetAdminPKBytes()
}

// GetOrCreateAccount return default account if not exist
func (db *mockDB) GetOrCreateAccount(address common.Address) (common.Account, error) {
	return db.db.GetOrCreateAccount(address)
}

// NewWriteBatch new a write batch
func (db *mockDB) NewWriteBatch() db.WriteBatch {
	wb := db.db.NewWriteBatch()
	return &mockWriteBatch{
		wb: wb,
	}
}

type mockWriteBatch struct {
	wb db.WriteBatch
}

func (wb *mockWriteBatch) AddBlock(block *core.Block) error {
	return wb.wb.AddBlock(block)
}
func (wb *mockWriteBatch) SetConsensusBlock(id string, num uint64) {
	wb.wb.SetConsensusBlock(id, num)
}
func (wb *mockWriteBatch) UpdateChannel(id string, profile *cc.Profile) error {
	time.Sleep(100 * time.Millisecond)
	return wb.wb.UpdateChannel(id, profile)
}
func (wb *mockWriteBatch) UpdateAccounts(accounts ...common.Account) error {
	return wb.wb.UpdateAccounts(accounts...)
}

func (wb *mockWriteBatch) SetAccount(account common.Account) error {
	return wb.wb.SetAccount(account)
}

func (wb *mockWriteBatch) SetAssetAdmin(pk crypto.PublicKey) error {
	return wb.wb.SetAssetAdmin(pk)
}
func (wb *mockWriteBatch) Put(key, value []byte) {
	wb.wb.Put(key, value)
}
func (wb *mockWriteBatch) Sync() error {
	return wb.wb.Sync()
}

func TestMockAddBlock(t *testing.T) {
	os.RemoveAll(".db")
	os.RemoveAll(".blocks")
	c, err := NewCoordinator(".db", &config.Config{
		BlockChain: config.BlockChainConfig{
			Path:         ".blocks",
			BatchTimeout: 10,
			BatchSize:    10,
		},
	}, &config.ConsensusConfig{
		Type: config.SOLO,
	})
	require.NoError(t, err)
	// update db
	c.db = &mockDB{
		db: c.db,
	}
	c.GM.db = c.db
	c.AM.db = c.db
	c.CM.db = c.db
	c.Consensus = new(mockConsensus)
	// magic because we want to make hub.Wait can be done
	c.GM.hub.Done("0b27ed8e359b7dc1b558cd4e87614180226244928185d95b1053cda2d4967712", nil)
	c.GM.hub.Done("4dee50f951867dc260d30311759310f3f0507e98534a082b0d1958f5fbd1a627", nil)
	c.GM.hub.Done("f96e6587ea3c2122ed59016dd52563a631a9c987f179fe7e8db75971a0bc8716", nil)
	c.Start()
	time.Sleep(1 * time.Second)
	os.RemoveAll(".db")
	os.RemoveAll(".blocks")
}
