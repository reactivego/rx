package rx

import (
	"sync"
	"sync/atomic"
)

// Subscriber is a Subscribable that allows construction of a Subscriber tree.
type Subscriber interface {
	// A Subscriber is a Subscribable
	Subscribable

	// Add will create and return a new child Subscriber setup in such a way that
	// calling Unsubscribe on the parent will also call Unsubscribe on the child.
	// Calling the Unsubscribe method on the child will NOT propagate to the
	// parent!
	Add() Subscriber

	// OnUnsubscribe will add the given callback function to the Subscriber.
	// The callback will be called when either the Unsubscribe of the parent
	// or of the Subscriber itself is called. If the subscription was already
	// canceled, then the callback function will just be called immediately.
	OnUnsubscribe(callback func())
}

const (
	subscribed = iota
	unsubscribed
)

type subscriber struct {
	state int32

	sync.Mutex
	callbacks []func()
}

func (s *subscriber) Subscribed() bool {
	return atomic.LoadInt32(&s.state) == subscribed
}

func (s *subscriber) Unsubscribe() {
	if atomic.CompareAndSwapInt32(&s.state, subscribed, unsubscribed) {
		s.Lock()
		for _, cb := range s.callbacks {
			cb()
		}
		s.callbacks = nil
		s.Unlock()
	}
}

func (s *subscriber) Canceled() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscriber) Add() Subscriber {
	child := &subscriber{}
	s.Lock()
	if atomic.LoadInt32(&s.state) != subscribed {
		child.Unsubscribe()
	} else {
		s.callbacks = append(s.callbacks, child.Unsubscribe)
	}
	s.Unlock()
	return child
}

func (s *subscriber) OnUnsubscribe(callback func()) {
	if callback == nil {
		return
	}
	s.Lock()
	if atomic.LoadInt32(&s.state) == subscribed {
		s.callbacks = append(s.callbacks, callback)
	} else {
		callback()
	}
	s.Unlock()
}
