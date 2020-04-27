package channel

import (
	"encoding/json"
	cc "madledger/blockchain/config"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"madledger/orderer/config"
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
	if num != 1 {
		time.Sleep(10 * time.Second)
	}
	time.Sleep(time.Duration(util.RandNum(1000)) * time.Millisecond)
	switch channelID {
	case core.GLOBALCHANNELID:
		// construct a special global block
		var txs = make([]*core.Tx, 0)
		// // first we should run config block
		// payload, _ := json.Marshal(core.GlobalTxPayload{
		// 	ChannelID: core.CONFIGCHANNELID,
		// 	Number:    1,
		// })
		// txs = append(txs, &core.Tx{
		// 	ID: "global1",
		// 	Data: core.TxData{
		// 		ChannelID: core.GLOBALCHANNELID,
		// 		Payload:   payload,
		// 	},
		// })
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
		return &mockBlock{
			num: 1,
			txs: txs,
		}, nil
	case core.CONFIGCHANNELID:
		time.Sleep(1 * time.Second)
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
		return &mockBlock{
			num: 1,
			txs: txs,
		}, nil
	case core.ASSETCHANNELID:
	default: //user channel
		// fmt.Println("here is ", channelID)
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
	c.Consensus = new(mockConsensus)
	// then we set test manager
	testManager, err := NewManager("test", c)
	require.NoError(t, err)
	c.setChannel(testManager.ID, testManager)
	// magic because we want to make hub.Wait can be done
	c.GM.hub.Done("0b27ed8e359b7dc1b558cd4e87614180226244928185d95b1053cda2d4967712", nil)
	c.GM.hub.Done("4dee50f951867dc260d30311759310f3f0507e98534a082b0d1958f5fbd1a627", nil)
	c.GM.hub.Done("f96e6587ea3c2122ed59016dd52563a631a9c987f179fe7e8db75971a0bc8716", nil)
	c.Start()
	time.Sleep(5 * time.Second)
	os.RemoveAll(".db")
	os.RemoveAll(".blocks")
	// t.Fail()
}
