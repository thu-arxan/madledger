// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto/hash"
	"madledger/common/event"
	"madledger/common/util"
	"sort"
	"sync"
)

// channel caches blocks, minBlock, and puts data into db
type channel struct {
	sync.RWMutex
	id        uint64
	db        *DB
	channelID string
	hub       *event.Hub

	minBlock uint64            // minBlock is the min block number that the blockchain system needed
	blocks   map[uint64]*Block // blockNum => Block
	blockCh  chan *Block
}

func newChannel(id uint64, channelID string, db *DB) *channel {
	ch := &channel{
		id:        id,
		db:        db,
		channelID: channelID,
		hub:       event.NewHub(),
		blocks:    make(map[uint64]*Block),
		blockCh:   make(chan *Block, 2048),
	}
	// todo: why here + 1?
	// to avoid add replicated Block when raft is not leader before closed
	// minBlock should be zero or chainNum + 1
	ch.minBlock = db.GetMinBlock(channelID)
	if ch.minBlock != 0 {
		ch.minBlock++
	}
	return ch
}

func (c *channel) addBlock(block *Block) {
	c.Lock()
	defer c.Unlock()

	// TODO: Should we not only use sm3?
	hash := string(hash.Hash(block.Bytes()))

	if util.Contain(c.blocks, block.GetNumber()) {
		c.hub.Done(hash, &event.Result{
			Err: errors.New("Duplicated block"),
		})
		return
	}

	c.blocks[block.Num] = block
	if block.GetNumber() >= c.minBlock {
		c.db.AddBlock(block)
	}
	c.hub.Done(hash, nil)
	// todo: cache here? safety?
	c.blockCh <- block
}

func (c *channel) getBlock(num uint64) *Block {
	c.RLock()
	defer c.RUnlock()
	return c.blocks[num]
}

func (c *channel) containsBlock(num uint64) bool {
	c.RLock()
	defer c.RUnlock()
	return util.Contain(c.blocks, num)
}

// gc delete block[num] in cache
func (c *channel) gc(num uint64) {
	c.Lock()
	defer c.Unlock()

	delete(c.blocks, num)
}

func (c *channel) getMinBlock() uint64 {
	c.RLock()
	defer c.RUnlock()
	return c.minBlock
}

func (c *channel) setMinBlock(num uint64) {
	c.Lock()
	defer c.Unlock()
	c.minBlock = num
	c.db.SetMinBlock(c.channelID, num)
}

func (c *channel) BlockCh() chan *Block {
	return c.blockCh
}

func (c *channel) marshal() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.Marshal(c.blocks)
}

func (c *channel) unmarshal(data []byte) error {
	var blocks map[uint64]*Block
	if err := json.Unmarshal(data, &blocks); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	for num, block := range blocks {
		if num >= c.minBlock && !util.Contain(c.blocks, num) {
			c.blocks[num] = block
		}
	}

	// todo: is it correct and efficient?
	// send recovered block:
	go func(c *channel) {
		var nums = make([]uint64, 0)
		for num := range blocks {
			nums = append(nums, num)
		}
		sort.Slice(nums, func(i, j int) bool {
			return nums[i] < nums[j]
		})

		for _, num := range nums {
			block := blocks[num]
			if block.GetNumber() >= c.getMinBlock() {
				c.blockCh <- block
			}
		}
	}(c)

	return nil
}

func (c *channel) notifyLater(block *Block) {
	c.blockCh <- block
}

// TODO: should we not only use sm3
func (c *channel) watch(block *Block) error {
	hash := string(hash.Hash(block.Bytes()))
	res := c.hub.Watch(hash, nil)
	if res == nil {
		return nil
	}
	return res.(*event.Result).Err
}

func (c *channel) fetchBlockDone(num uint64) {
	c.Lock()
	defer c.Unlock()

	// todo: check num > current minBlock?
	c.minBlock = num + 1
	c.db.SetMinBlock(c.channelID, num+1)
	delete(c.blocks, num)
}
