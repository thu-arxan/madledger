package performance

import (
	"fmt"
	"log"
	"madledger/core"
	"madledger/tests/performance/bft"
	"os"
	"sync"
	"testing"
	"time"

	"madledger/tests/performance/raft"
	"madledger/tests/performance/solo"

	cutil "madledger/client/util"

	"net/http"
	_ "net/http/pprof"

	"github.com/stretchr/testify/require"
)

var (
	consensus   = "solo"
	peerNum     = 3
	channelSize = 10
	clientSize  = 200
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func TestInit(t *testing.T) {
	os.Remove(logPath)
	switch consensus {
	case "solo":
		require.NoError(t, solo.Init(clientSize))
		require.NoError(t, solo.StartOrderers())
		require.NoError(t, solo.StartPeers())
	case "raft":
		require.NoError(t, raft.Init(clientSize, peerNum))
		require.NoError(t, raft.StartOrderers())
		require.NoError(t, raft.StartPeers(peerNum))
	case "bft":
		require.NoError(t, bft.Init(clientSize, peerNum))
		require.NoError(t, bft.StartOrderers())
		require.NoError(t, bft.StartPeers(peerNum))
	default:
		panic("Unsupport consensus")
	}

}

func TestCreateChannel(t *testing.T) {
	var clients = getClients()
	for i := 0; i < channelSize; i++ {
		require.NoError(t, clients[0].CreateChannel(fmt.Sprintf("test%d", i), true, nil, nil))
	}
}

func TestCreateContract(t *testing.T) {
	var clients = getClients()
	for i := 0; i < channelSize; i++ {
		CreateContract(t, fmt.Sprintf("test%d", i), clients[0])
	}
}

func TestPerformance(t *testing.T) {
	var wg sync.WaitGroup
	var callSize = 40
	clients := getClients()
	var txs = make([][]*core.Tx, clientSize)
	// create txs
	for i := range txs {
		txs[i] = CreateCallContractTx(fmt.Sprintf("test%d", i%channelSize), clients[i], callSize)
	}
	begin := time.Now()
	for i := range clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := clients[i]
			AddTxs(t, client, txs[i])
		}(i)
	}
	wg.Wait()
	duration := (int64)(time.Since(begin)) / 1e6
	tps := int64(callSize*len(clients)*1e3) / (duration)
	table := cutil.NewTable()
	table.SetHeader("Consensus", "Size", "Time", "TPS")
	table.AddRow(consensus, callSize*len(clients), fmt.Sprintf("%v", time.Since(begin)), fmt.Sprintf("%d", tps))
	require.NoError(t, writeLog(table.ToString()))
}

func TestEnd(t *testing.T) {
	switch consensus {
	case "solo":
		solo.Clean()
	case "raft":
		raft.Clean()
	case "bft":
		bft.Clean()
	default:
		panic("Unsupported consensus")
	}
}
