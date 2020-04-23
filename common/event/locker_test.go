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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLockerWithoutRelay(t *testing.T) {
	l := NewLocker()
	l.Unlock("subject", 10)
	require.Len(t, l.Wait("subject", 0), 0)
	l.Wait("subject", 8)
	l.Wait("subject", 10)
	var wg sync.WaitGroup
	// should be ok
	wg.Add(2)
	go func() {
		defer wg.Done()
		l.Wait("subject", 13)
	}()
	go func() {
		defer wg.Done()
		l.Unlock("subject", 15)
	}()
	wg.Wait()
	// success one, and fail one
	var success = make(chan bool, 1)
	var fail = make(chan bool, 1)
	var count = 0
	go func() {
		l.Wait("another", 10)
		success <- true
	}()
	go func() {
		l.Wait("another", 15)
		fail <- true
	}()
	go func() {
		l.Unlock("another", 12)
	}()
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			require.EqualValues(t, 1, count)
			return
		case <-success:
			count = 1
		case <-fail:
			count = -1
		}
	}
}

func TestLockerWithRelay(t *testing.T) {
	l := NewLocker()
	l.Unlock("subject", 5, NewSubject("another", 10))
	subjects := l.Wait("subject", 1)
	require.Len(t, subjects, 1)
	require.EqualValues(t, "another", subjects[0].K)
	subjects = l.Wait("subject", 2)
	require.Len(t, subjects, 0)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		require.Len(t, l.Wait("subject", 10), 1)
	}()
	go func() {
		defer wg.Done()
		require.Len(t, l.Wait("subject", 11), 0)
	}()
	time.Sleep(100 * time.Millisecond)
	l.Unlock("subject", 15, NewSubject("another", 20))
	wg.Wait()
	// then we unlock while wait is not done
	l.Unlock("subject", 30, NewSubject("another", 25))
	require.Len(t, l.Wait("subject", 30), 1)
	// todo: we need more works
}

func TestBinaryInsert(t *testing.T) {
	// 1 -> 1,1
	var waits = binaryInsert(newWaitList(1), &waitChan{num: 1})
	require.Len(t, waits, 2)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 1)
	// 1 -> 0,1
	waits = binaryInsert(newWaitList(1), &waitChan{num: 0})
	require.Len(t, waits, 2)
	require.EqualValues(t, waits[0].num, 0)
	require.EqualValues(t, waits[1].num, 1)
	// 1 -> 1,2
	waits = binaryInsert(newWaitList(1), &waitChan{num: 2})
	require.Len(t, waits, 2)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 2)
	// 1,1 -> 1,1,2
	waits = binaryInsert(newWaitList(1, 1), &waitChan{num: 2})
	require.Len(t, waits, 3)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 1)
	require.EqualValues(t, waits[2].num, 2)
	// 1,1 -> 1,1,1
	waits = binaryInsert(newWaitList(1, 1), &waitChan{num: 1})
	require.Len(t, waits, 3)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 1)
	require.EqualValues(t, waits[2].num, 1)
	// 1,1 -> 0,1,1
	waits = binaryInsert(newWaitList(1, 1), &waitChan{num: 0})
	require.Len(t, waits, 3)
	require.EqualValues(t, waits[0].num, 0)
	require.EqualValues(t, waits[1].num, 1)
	require.EqualValues(t, waits[2].num, 1)
	// 1,3 -> 1,2,3
	waits = binaryInsert(newWaitList(1, 3), &waitChan{num: 2})
	require.Len(t, waits, 3)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 2)
	require.EqualValues(t, waits[2].num, 3)
	// 1,3,3,5 -> 1,3,3,4,5
	waits = binaryInsert(newWaitList(1, 3, 3, 5), &waitChan{num: 4})
	require.Len(t, waits, 5)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 3)
	require.EqualValues(t, waits[2].num, 3)
	require.EqualValues(t, waits[3].num, 4)
	require.EqualValues(t, waits[4].num, 5)
	// 1,3,3,5 -> 1,3,3,5,7
	waits = binaryInsert(newWaitList(1, 3, 3, 5), &waitChan{num: 7})
	require.Len(t, waits, 5)
	require.EqualValues(t, waits[0].num, 1)
	require.EqualValues(t, waits[1].num, 3)
	require.EqualValues(t, waits[2].num, 3)
	require.EqualValues(t, waits[3].num, 5)
	require.EqualValues(t, waits[4].num, 7)
}

func newWaitList(nums ...uint64) []*waitChan {
	var waits = make([]*waitChan, 0)
	for i := range nums {
		waits = append(waits, &waitChan{
			num: nums[i],
		})
	}
	return waits
}
