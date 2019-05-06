// +build js,wasm

package subscriber

import (
	"sync/atomic"
)

const (
	subscribed = iota
	unsubscribed
)

type subscription struct {
	state int32

	callbacks []func()
}

func (s *subscription) Unsubscribe() {
	if atomic.CompareAndSwapInt32(&s.state, subscribed, unsubscribed) {
		for _, cb := range s.callbacks {
			cb()
		}
		s.callbacks = nil
	}
}

func (s *subscription) Closed() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscription) Wait() {
	// Not implemented 
}

func (s *subscription) Add(callback func()) Subscriber {
	child := &subscription{callbacks: []func(){callback}}
	if s.Closed() {
		child.Unsubscribe()
	} else {
		s.callbacks = append(s.callbacks, child.Unsubscribe)
	}
	return child
}
