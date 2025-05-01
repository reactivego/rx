package rx

import (
	"sync"
	"sync/atomic"
)

// Subscriber is a subscribable entity that allows construction of a Subscriber tree.
type Subscriber interface {
	// Subscribed returns true if the subscriber is in a subscribed state.
	// Returns false once Unsubscribe has been called.
	Subscribed() bool

	// Unsubscribe changes the state to unsubscribed and executes all registered
	// callback functions. Does nothing if already unsubscribed.
	Unsubscribe()

	// Add creates and returns a new child Subscriber.
	// If the parent is already unsubscribed, the child will be created in an
	// unsubscribed state. Otherwise, the child will be unsubscribed when the parent
	// is unsubscribed.
	Add() Subscriber

	// OnUnsubscribe registers a callback function to be executed when Unsubscribe is called.
	// If the subscriber is already unsubscribed, the callback is executed immediately.
	// If callback is nil, this method does nothing.
	OnUnsubscribe(callback func())
}

const (
	subscribed = iota
	unsubscribed
)

type subscriber struct {
	state atomic.Int32
	sync.Mutex
	callbacks []func()
}

func (s *subscriber) Subscribed() bool {
	return s.state.Load() == subscribed
}

func (s *subscriber) Unsubscribe() {
	if s.state.CompareAndSwap(subscribed, unsubscribed) {
		s.Lock()
		for _, cb := range s.callbacks {
			cb()
		}
		s.callbacks = nil
		s.Unlock()
	}
}

func (s *subscriber) Add() Subscriber {
	child := &subscriber{}
	s.Lock()
	if s.state.Load() != subscribed {
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
	if s.state.Load() == subscribed {
		s.callbacks = append(s.callbacks, callback)
	} else {
		callback()
	}
	s.Unlock()
}
