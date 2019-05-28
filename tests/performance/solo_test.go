package performance

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"madledger/tests/performance/solo"

	cutil "madledger/client/util"

	"github.com/stretchr/testify/require"
)

func TestSoloInit(t *testing.T) {
	os.Remove(logPath)
	require.NoError(t, solo.Init())
	require.NoError(t, solo.StartOrderers())
	require.NoError(t, solo.StartPeers())
}

func TestCreateChannel(t *testing.T) {
	clients := solo.GetClients()
	require.NoError(t, clients[0].CreateChannel("test", true, nil, nil))
}

func TestCreateContract(t *testing.T) {
	clients := solo.GetClients()
	CreateContract(t, "test", clients[0])
}

func TestPerformance(t *testing.T) {
	var wg sync.WaitGroup
	var callSize = 40
	clients := solo.GetClients()
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
	table.AddRow("Solo", callSize*len(clients), fmt.Sprintf("%v", time.Since(begin)), fmt.Sprintf("%d", tps))
	require.NoError(t, writeLog(table.ToString()))
}

func TestSoloEnd(t *testing.T) {
	solo.Clean()
}
