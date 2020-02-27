package event

import (
	"fmt"
	"madledger/common/util"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	var hub = NewHub()
	var id = util.RandomString(10)
	go func() {
		time.Sleep(100 * time.Millisecond)
		hub.Done(id, 1)
	}()
	var finish = make(chan bool, 1)
	go func() {
		num := hub.Watch(id, nil).(int)
		if num != 1 {
			t.Fatal()
		}
		finish <- true
	}()
	num := hub.Watch(id, nil).(int)
	if num != 1 {
		t.Fatal()
	}
	<-finish
	hub.Done(id, 2)
	num = hub.Watch(id, nil).(int)
	if num != 1 {
		t.Fatal()
	}
}

func TestRegister(t *testing.T) {
	var hub = NewHub()
	var topic = util.RandomString(10)
	var finish = make(chan bool, 1)
	go func() {
		ch, token := hub.Register(topic)
		for {
			select {
			case msg := <-ch:
				s := msg.(string)
				if s == "5" {
					hub.UnRegister(topic, token)
					finish <- true
					return
				}
			}
		}
	}()
	go func() {
		for i := 0; i <= 5; i++ {
			time.Sleep(100 * time.Millisecond)
			hub.Broadcast(topic, fmt.Sprintf("%d", i))
		}
	}()
	<-finish
}
