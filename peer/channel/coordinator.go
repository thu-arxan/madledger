package channel

import (
	"madledger/common/util"
	"sync"
)

// Coordinator is responsible for coordinate.
type Coordinator struct {
	lock   sync.Mutex
	states map[string]*State
}

// StateCode represent the code of state
type StateCode int

// All States
const (
	Waitting StateCode = iota
	Runable
)

// Dependency defines the channel and block that depends on
type Dependency struct {
	ChannelID string
	Num       uint64
}

// State represents the state of channel
type State struct {
	num  uint64
	code StateCode
	// hashes is not working now
	hashes map[uint64][]byte
	// dependencies is not working now
	dependencies []Dependency
}

// NewCoordinator is the constructor of Coordinator
func NewCoordinator() *Coordinator {
	c := new(Coordinator)
	c.states = make(map[string]*State)
	return c
}

// CanRun return runable
func (c *Coordinator) CanRun(channelID string, num uint64) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if util.Contain(c.states, channelID) {
		state := c.states[channelID]
		if num < state.num {
			return true
		}
		if num > state.num {
			return false
		}
		return state.code == Runable
	}
	return false
}

// Locks will lock some channels because all blocks should run after the config
// blocks are all done.
func (c *Coordinator) Locks() {
	c.lock.Lock()
	defer c.lock.Unlock()

}

// Unlocks will unlock some channels
func (c *Coordinator) Unlocks(nums map[string]uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for channel, num := range nums {
		if util.Contain(c.states, channel) {
			state := c.states[channel]
			if num >= state.num {
				state.num = num
				state.code = Runable
			}
		} else {
			state := new(State)
			state.num = num
			state.code = Runable
			c.states[channel] = state
		}
	}
}
