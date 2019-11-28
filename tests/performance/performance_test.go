package performance

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"madledger/tests/performance/raft"
	"madledger/tests/performance/solo"

	client "madledger/client/lib"
	cutil "madledger/client/util"

	"github.com/stretchr/testify/require"
)

var consensus = "raft"

func TestInit(t *testing.T) {
	os.Remove(logPath)
	switch consensus {
	case "solo":
		require.NoError(t, solo.Init())
		require.NoError(t, solo.StartOrderers())
		require.NoError(t, solo.StartPeers())
	case "raft":
		require.NoError(t, raft.Init())
		require.NoError(t, raft.StartOrderers())
		require.NoError(t, raft.StartPeers())
	case "bft":
		panic("Unsupport now")
	default:
		panic("Unsupport consensus")
	}

}

func TestCreateChannel(t *testing.T) {
	var clients = getClients()
	require.NoError(t, clients[0].CreateChannel("test", true, nil, nil))
}

func TestCreateContract(t *testing.T) {
	var clients = getClients()
	CreateContract(t, "test", clients[0])
}

func TestPerformance(t *testing.T) {
	var wg sync.WaitGroup
	var callSize = 40
	clients := getClients()
	begin := time.Now()
	for i := range clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := clients[i]
			CallContract(t, "test", client, callSize)
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

func TestSoloEnd(t *testing.T) {
	switch consensus {
	case "solo":
		solo.Clean()
	case "raft":
		raft.Clean()
	case "bft":
	default:
		panic("Unsupport consensus")
	}
}

func getClients() []*client.Client {
	var clients []*client.Client
	switch consensus {
	case "solo":
		clients = solo.GetClients()
	case "raft":
		clients = raft.GetClients()
	case "bft":
	default:
		panic("Unsupport consensus")
	}
	return clients
}
