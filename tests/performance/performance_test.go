package performance

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	client "madledger/client/lib"
	cutil "madledger/client/util"
	"madledger/common/util"

	"github.com/stretchr/testify/require"
)

func TestInitCircumstance(t *testing.T) {
	os.Remove(logPath)
	err := initDir(".orderer")
	require.NoError(t, err)
	err = initDir(".peer")
	require.NoError(t, err)
	// then start necessary orderer and peer
	err = startSoloOrderer()
	require.NoError(t, err)
	err = startSoloPeer()
	require.NoError(t, err)
	for i := range clients {
		cfgPath, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/performance/.clients/%d/client.yaml", i), gopath)
		c, err := client.NewClient(cfgPath)
		require.NoError(t, err)
		clients[i] = c
	}
}

func TestCreateChannel(t *testing.T) {
	require.NoError(t, clients[0].CreateChannel("test", true, nil, nil))
}

func TestCreateContract(t *testing.T) {
	createContract(t, "test", clients[0])
}

func TestPerformance(t *testing.T) {
	var wg sync.WaitGroup
	var callSize = 40
	begin := time.Now()
	for i := range clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := clients[i]
			callContract(t, "test", client, callSize)
		}(i)
	}
	wg.Wait()
	duration := (int64)(time.Since(begin)) / 1e6
	tps := int64(callSize*len(clients)*1e3) / (duration)
	table := cutil.NewTable()
	table.SetHeader("Size", "Time", "TPS")
	table.AddRow(callSize*len(clients), fmt.Sprintf("%v", time.Since(begin)), fmt.Sprintf("%d", tps))
	require.NoError(t, writeLog(table.ToString()))
}

func TestRemoveData(t *testing.T) {
	os.RemoveAll(".clients")
	os.RemoveAll(".orderer")
	os.RemoveAll(".peer")
}
