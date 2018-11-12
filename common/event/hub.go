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
// todo: gc is still not done yet.
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

// Watches will watch many events and return when all events are finished.
func (h *Hub) Watches(ids []string) *Result {
	h.lock.Lock()

	// first check if there are event finished and contain error
	for _, id := range ids {
		// finished and err is not nil and err.Error() is not empty
		if util.Contain(h.finish, id) && h.finish[id].Err != nil && h.finish[id].Err.Error() != "" {
			defer h.lock.Unlock()
			return h.finish[id]
		}
	}
	// then ignore events that are finished and finish right
	var unfinish []string
	for _, id := range ids {
		if !util.Contain(h.finish, id) {
			unfinish = append(unfinish, id)
			if !util.Contain(h.events, id) {
				h.events[id] = make([]chan bool, 0)
			}
		}
	}

	// then set chans to watch these unfinished events
	var chs = make(map[string]chan bool)
	for _, id := range unfinish {
		ch := make(chan bool, 1)
		h.events[id] = append(h.events[id], ch)
		chs[id] = ch
	}
	h.lock.Unlock()

	// waitting events finish
	var errs = make(chan error, len(unfinish))
	for id := range chs {
		go func(id string) {
			<-chs[id]
			h.lock.Lock()
			result := h.finish[id]
			h.lock.Unlock()
			if result.Err != nil && result.Err.Error() != "" {
				errs <- result.Err
			} else {
				errs <- nil
			}
		}(id)
	}

	var count = 0
	var result = Result{
		Err: nil,
	}
	for {
		err := <-errs
		if err != nil {
			result.Err = err
			break
		}
		count++
		if count == len(unfinish) {
			break
		}
	}

	return &result
}
