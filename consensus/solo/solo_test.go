package solo

import (
	"madledger/consensus"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	sc consensus.Consensus
)

func TestStart(t *testing.T) {
	var err error
	cfg := make(map[string]consensus.Config)
	cfg["test"] = consensus.DefaultConfig()
	sc, err = NewConsensus(cfg)
	require.NoError(t, err)
	sc.Start()
}

func TestAddTx(t *testing.T) {
	var txSize = 1024
	var txs [][]byte
	var success = make(map[string]int)
	var lock sync.Mutex

	for i := 0; i < txSize; i++ {
		tx := randomTx()
		success[string(tx)] = 0
		txs = append(txs, tx)
	}

	var wg sync.WaitGroup
	// each tx send 3 times
	for i := range txs {
		for m := 0; m < 3; m++ {
			wg.Add(1)
			tx := txs[i]
			go func() {
				defer wg.Done()
				if err := sc.AddTx("test", tx); err == nil {
					lock.Lock()
					success[string(tx)]++
					lock.Unlock()
				}
			}()
		}
	}
	wg.Wait()

	for i := range success {
		require.Equal(t, 1, success[i])
	}
}

func TestEnd(t *testing.T) {
	err := sc.Stop()
	require.NoError(t, err)
}

func randomTx() []byte {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 32)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return []byte(string(b))
}
