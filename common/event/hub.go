// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.
package event

import (
	"sync"
)

// Hub provide an easy way to register event
type Hub struct {
	lock *sync.Mutex
	// watch & done
	events map[string][]chan interface{}
	finish map[string]interface{}
	// register & broadcast
	tally  int
	topics map[string][]int
	chs    map[int]chan interface{}
}

// CallBack is callback function
type CallBack func()

// NewHub is the constructor of Hub
func NewHub() *Hub {
	return &Hub{
		lock:   new(sync.Mutex),
		events: make(map[string][]chan interface{}),
		finish: make(map[string]interface{}),
		topics: make(map[string][]int),
		chs:    make(map[int]chan interface{}),
	}
}

// Done done an event
// Note: Done same thing twice would be ignored.
func (h *Hub) Done(id string, res interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.finish[id]; ok {
		return
	}
	h.finish[id] = res

	if chs, ok := h.events[id]; ok {
		for _, ch := range chs {
			ch <- res
		}
	}

	delete(h.events, id)
}

// Watch watch an event
// Note: CallBack function is not called only after watch done but also succeed register watch event,
// and it should be setted carefully
func (h *Hub) Watch(id string, rec CallBack) interface{} {
	h.lock.Lock()

	if _, ok := h.finish[id]; ok {
		if rec != nil {
			rec()
		}
		defer h.lock.Unlock()
		return h.finish[id]
	}

	if _, ok := h.events[id]; !ok {
		h.events[id] = make([]chan interface{}, 0)
	}

	ch := make(chan interface{}, 1)
	h.events[id] = append(h.events[id], ch)

	if rec != nil {
		rec()
	}
	h.lock.Unlock()

	result := <-ch
	return result
}

// Broadcast boradcast msg of topics
func (h *Hub) Broadcast(topic string, msg interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if tokens, ok := h.topics[topic]; ok {
		for i := range tokens {
			// TODO: if the receiver refused to receive the msg, we should deal this
			if ch, ok := h.chs[tokens[i]]; ok {
				ch <- msg
			}
		}
	}
}

// Register will register an id, the register should hold the token to delete itself
func (h *Hub) Register(topic string) (ch chan interface{}, token int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.tally++

	token = h.tally
	ch = make(chan interface{}, 128)
	if _, ok := h.topics[topic]; !ok {
		h.topics[topic] = make([]int, 0)
	}
	h.topics[topic] = append(h.topics[topic], h.tally)
	h.chs[h.tally] = ch

	return
}

// UnRegister unregister token
func (h *Hub) UnRegister(topic string, token int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if s, ok := h.topics[topic]; ok {
		s = removeFromSlice(s, token)
		if len(s) == 0 {
			delete(h.topics, topic)
		} else {
			h.topics[topic] = s
		}
	}
	if _, ok := h.chs[token]; ok {
		delete(h.chs, token)
	}
}

func removeFromSlice(s []int, elem int) []int {
	var result = make([]int, 0)
	for i := range s {
		if s[i] != elem {
			result = append(result, s[i])
		}
	}
	return result
}
