package event

import (
	"fmt"
	"madledger/common/util"
	"sync"
)

// Hub manage all events
type Hub struct {
	lock   *sync.Mutex
	events map[string][]chan bool
	finish map[string]*Result
}

// NewHub is the constructor of Hub
func NewHub() *Hub {
	return &Hub{
		lock:   new(sync.Mutex),
		finish: make(map[string]*Result),
		events: make(map[string][]chan bool),
	}
}

// Done set the result of event.
// One event could only set done once now.
func (h *Hub) Done(id string, result *Result) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// event is not finished and
	if !util.Contain(h.finish, id) {
		h.finish[id] = result
		if result == nil {
			h.finish[id] = NewResult(nil)
		}

		for _, ch := range h.events[id] {
			ch <- true
		}
	}
}

// Watch will watch an event.
// gc is still now be done yet.
func (h *Hub) Watch(id string, wc *WatchConfig) *Result {
	h.lock.Lock()

	if wc == nil {
		wc = DefaultWatchConfig()
	}

	if util.Contain(h.finish, id) {
		defer h.lock.Unlock()
		return h.finish[id]
	}

	if !util.Contain(h.events, id) {
		h.events[id] = make([]chan bool, 0)
	}

	if wc.Single && len(h.events) != 0 {
		defer h.lock.Unlock()
		return NewResult(fmt.Errorf("Duplicate watch is not allowed in single mode"))
	}

	ch := make(chan bool, 1)
	h.events[id] = append(h.events[id], ch)
	h.lock.Unlock()

	<-ch
	h.lock.Lock()
	result := h.finish[id]
	defer h.lock.Unlock()

	return result
}
