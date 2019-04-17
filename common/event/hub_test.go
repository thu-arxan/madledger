package event

import (
	"fmt"
	"madledger/common/util"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	eventSize   = 2048
	eventLength = 32
)

func TestWatch(t *testing.T) {
	hub := NewHub()
	events := make([]string, eventSize)
	// initial events
	for i := range events {
		events[i] = util.RandomString(eventLength)
	}

	var wg sync.WaitGroup
	// register events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < eventSize*10; i++ {
			event := events[i%eventSize]
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result := hub.Watch(event, nil)
				require.EqualError(t, result.Err, fmt.Sprintf("Error is %d", i%eventSize))
			}(i)
		}
	}()

	// finish events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < eventSize; i++ {
			event := events[i]
			wg.Add(1)
			go func(i int) {
				wg.Done()
				hub.Done(event, &Result{
					Err: fmt.Errorf("Error is %d", i),
				})
			}(i)
		}
	}()

	wg.Wait()
	// test if all events is clean
	for event := range hub.events {
		require.Len(t, hub.events[event], 0)
	}
}

func TestWatches(t *testing.T) {
	hub := NewHub()
	events := make([]string, eventSize)
	// initial events
	for i := range events {
		events[i] = util.RandomString(eventLength)
	}

	var wg sync.WaitGroup
	// register events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < eventSize*10; i++ {
			event := events[i%eventSize]
			var es []string
			es = append(es, events[i%eventSize])
			// only watch succeed event if i%2=0
			if i%2 == 0 {
				es = append(es, events[(i+2)%eventSize])
			} else {
				es = append(es, events[(i+1)%eventSize])
			}
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				randomSleep()
				result := hub.Watch(event, nil)
				if i%2 == 0 {
					require.NoError(t, result.Err)
				} else {
					require.EqualError(t, result.Err, fmt.Sprintf("Error is %d", i%eventSize))
				}
			}(i)
		}
	}()

	// finish events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < eventSize; i++ {
			event := events[i]
			wg.Add(1)
			go func(i int) {
				wg.Done()
				randomSleep()
				if i%2 == 0 {
					hub.Done(event, &Result{
						Err: nil,
					})
				} else {
					hub.Done(event, &Result{
						Err: fmt.Errorf("Error is %d", i),
					})
				}
			}(i)
		}
	}()

	wg.Wait()
	// test if all events is clean
	for event := range hub.events {
		require.Len(t, hub.events[event], 0)
	}
}

func randomSleep() {
	time.Sleep(time.Duration(util.RandNum(500)) * time.Millisecond)
}
