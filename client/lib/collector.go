package lib

import (
	"encoding/json"
	"errors"
	"madledger/common/util"
	"sync"
)

// Collector means to collect results.
// Once it collects enough results then it will generate final results.
type Collector struct {
	lock    sync.Mutex
	result  interface{}
	results map[string]int
	errors  []error
	max     int
	finish  bool
	// finish channel
	fc chan bool
}

// NewCollector is the constructor of Collector
func NewCollector(max int) *Collector {
	collector := new(Collector)
	collector.max = max
	collector.finish = false
	collector.results = make(map[string]int)
	collector.fc = make(chan bool, 1)

	return collector
}

// Add add result
func (c *Collector) Add(result interface{}, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.finish {
		return
	}
	if err != nil {
		c.errors = append(c.errors, err)
	} else {
		data, _ := json.Marshal(result)
		s := util.Hex(data)
		if util.Contain(c.results, s) {
			c.results[s]++
		} else {
			c.results[s] = 1
		}
		if c.results[s] >= (c.max/2 + 1) {
			c.result = result
			c.finish = true
			c.fc <- true
		}
	}
	if len(c.errors) >= (c.max/2+1) || (len(c.results)+len(c.errors)) >= c.max {
		c.finish = true
		c.fc <- true
	}
}

// Wait wait the final result is decided
func (c *Collector) Wait() (interface{}, error) {
	<-c.fc
	if c.result != nil {
		return c.result, nil
	}
	if len(c.errors) >= (c.max/2 + 1) {
		return nil, c.errors[0]
	}
	return nil, errors.New("Failed to get enough same results")
}
