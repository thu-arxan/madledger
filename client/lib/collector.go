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

	total int
	min   int

	finish bool
	// finish channel
	fc chan bool
}

// NewCollector is the constructor of Collector
// total means total result we may collect, min means if we has min same result then we finish the collector
// Note: we will set min as total/2 + 1 if min <= 0
func NewCollector(total, min int) *Collector {
	collector := new(Collector)

	collector.results = make(map[string]int)

	collector.total = total
	if min <= 0 {
		min = total/2 + 1
	}
	collector.min = min

	collector.finish = false
	collector.fc = make(chan bool, 1)

	return collector
}

// Add add result
func (c *Collector) Add(result interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.finish {
		return
	}

	data, _ := json.Marshal(result)
	s := util.Hex(data)
	if util.Contain(c.results, s) {
		c.results[s]++
	} else {
		c.results[s] = 1
	}
	if c.results[s] >= c.min {
		c.result = result
		c.done()
	}

	var total int
	for _, num := range c.results {
		total += num
	}
	if len(c.errors)+total >= c.total {
		c.done()
	}
}

// AddError add an error
func (c *Collector) AddError(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.finish {
		return
	}

	c.errors = append(c.errors, err)
	if len(c.errors) >= c.min {
		c.done()
	}
}

// done set finish as true and add signal in fc
func (c *Collector) done() {
	c.finish = true
	c.fc <- true
}

// Wait wait the final result is decided
func (c *Collector) Wait() (interface{}, error) {
	<-c.fc
	if c.result != nil {
		return c.result, nil
	}
	if len(c.errors) >= c.min {
		return nil, c.errors[0]
	}
	return nil, errors.New("Failed to get enough same results")
}
