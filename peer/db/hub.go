package db

import (
	"madledger/common/util"
	"sync"
)

// Hub provide an easy way to register event
type Hub struct {
	lock   *sync.Mutex
	events map[string][]chan *TxStatus
	finish map[string]*TxStatus
}

// CallBack is callback function
type CallBack func()

// NewHub is the constructor of Hub
func NewHub() *Hub {
	return &Hub{
		lock:   new(sync.Mutex),
		events: make(map[string][]chan *TxStatus),
		finish: make(map[string]*TxStatus),
	}
}

// Done done an event
func (h *Hub) Done(id string, res *TxStatus) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.finish[id] = res

	if !util.Contain(h.events, id) {
		return
	}

	for _, ch := range h.events[id] {
		ch <- res
	}
}

// Watch watch an event
func (h *Hub) Watch(id string, rec CallBack) *TxStatus {
	h.lock.Lock()

	if util.Contain(h.finish, id) {
		if rec != nil {
			rec()
		}
		defer h.lock.Unlock()
		return h.finish[id]
	}

	if !util.Contain(h.events, id) {
		h.events[id] = make([]chan *TxStatus, 0)
	}

	ch := make(chan *TxStatus, 1)
	h.events[id] = append(h.events[id], ch)

	if rec != nil {
		rec()
	}
	h.lock.Unlock()

	result := <-ch
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.events[id]; ok {
		delete(h.events, id)
	}
	return result
}
