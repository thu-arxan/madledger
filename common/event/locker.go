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

// Locker will lock some thing until some things happens
type Locker struct {
	lock     sync.Mutex
	subjects map[string]uint64
	relays   map[string][]*Subject
	waits    map[string][]*waitChan
}

// Subject is the subject of locker content
type Subject struct {
	K string
	V uint64
}

// NewSubject is the constructor of Subject
func NewSubject(subject string, num uint64) *Subject {
	return &Subject{
		K: subject,
		V: num,
	}
}

type waitChan struct {
	num uint64
	ch  chan []*Subject
}

func newWaitChan(num uint64) *waitChan {
	return &waitChan{
		num: num,
		ch:  make(chan []*Subject, 1),
	}
}

// NewLocker is the constructor of Locker
func NewLocker() *Locker {
	return &Locker{
		subjects: make(map[string]uint64),
		relays:   make(map[string][]*Subject),
		waits:    make(map[string][]*waitChan),
	}
}

// Wait wait until somethings happens
func (l *Locker) Wait(subject string, num uint64) []*Subject {
	l.lock.Lock()
	if value, exist := l.subjects[subject]; exist {
		if value >= num {
			if subjects, exist := l.relays[subject]; exist {
				delete(l.relays, subject)
				l.lock.Unlock()
				return subjects
			}
			l.lock.Unlock()
			return nil
		}
	}
	// then register waitChan
	var waits = l.waits[subject]
	if len(waits) == 0 {
		waits = make([]*waitChan, 0)
	}
	waitCh := newWaitChan(num)
	waits = binaryInsert(waits, waitCh)
	l.waits[subject] = waits
	l.lock.Unlock()
	// wait signal
	subjects := <-waitCh.ch
	// close(waitCh.ch)
	if len(subjects) == 0 {
		return nil
	}
	return subjects
}

// Unlock unlock something, and it could relay subjects to the observer
// Note: The first observer or the minist num observer will relay subject, so if you want to relay subject then
// you should not wait subject with same num.
func (l *Locker) Unlock(subject string, num uint64, subjects ...*Subject) {
	l.lock.Lock()
	if value, exist := l.subjects[subject]; exist {
		// do nothing
		if value >= num {
			l.lock.Unlock()
			return
		}
	}
	l.subjects[subject] = num
	waits := l.waits[subject]
	if len(waits) != 0 {
		var newWaits = make([]*waitChan, 0)
		for i := range waits {
			if waits[i].num <= num {
				if i == 0 && len(subjects) != 0 {
					waits[i].ch <- subjects
				} else {
					waits[i].ch <- nil
				}
			} else {
				newWaits = append(newWaits, waits[i])
			}
		}
		if len(newWaits) == 0 {
			delete(l.waits, subject)
		} else {
			l.waits[subject] = newWaits
			if len(newWaits) == len(waits) { // subjects were not sent
				// todo: it only works when the user wait increment
				// todo: we need a better way to do this
				l.relays[subject] = subjects
			}
		}
	} else if len(subjects) != 0 {
		l.relays[subject] = subjects
	}
	// then unlock
	l.lock.Unlock()
}

func binaryInsert(waits []*waitChan, wait *waitChan) []*waitChan {
	if len(waits) == 0 {
		waits = append(waits, wait)
	}
	var begin = 0
	var end = len(waits) - 1
	var idx = begin
	for begin <= end {
		if waits[begin].num > wait.num {
			idx = begin
			break
		} else if waits[end].num <= wait.num {
			idx = end + 1
			break
		}
		mid := (begin + end) / 2
		if waits[mid].num <= wait.num {
			begin++
		} else {
			end--
		}
		idx = begin
	}
	waits = append(waits, nil)
	for i := len(waits) - 1; i > idx; i-- {
		waits[i] = waits[i-1]
	}
	waits[idx] = wait
	return waits
}
