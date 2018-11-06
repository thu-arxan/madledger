package event

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	eventSize = 2048
)

func TestHub(t *testing.T) {
	hub := NewHub()
	events := make([]string, eventSize)
	// initial events
	for i := range events {
		events[i] = randomStr()
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
				result := hub.Watch(event)
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
}

func randomStr() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 32)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
