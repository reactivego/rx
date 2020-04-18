package subscriber

import (
	"sync"
	"sync/atomic"
)

// Subscription is an interface that allows code to monitor and control a
// subscription it received.
type Subscription interface {
	// Unsubscribe will cancel the subscription (when one is active).
	// Subsequently it will then call Unsubscribe on all child subscriptions
	// added through Add. After a call to Unsubscribe returns, calling
	// Closed on the same interface or any of its child subscriptions
	// will return true. Unsubscribe can be safely called on a closed
	// subscription and performs no operation.
	Unsubscribe()

	// Closed returns true when the subscription has been canceled. There
	// is not need to check if the subscription is canceled before calling
	// Unsubscribe. This is an alias for the Canceled() method.
	Closed() bool

	// Canceled returns true when the subscription has been canceled. There
	// is not need to check if the subscription is canceled before calling
	// Unsubscribe.
	Canceled() bool

	// Wait will block the calling goroutine and wait for the Unsubscribe
	// method to be called on this subscription. Calling Wait on a
	// subscription that has already been closed will return immediately.
	Wait()
}

// Subscriber embeds a Subscription interface. Additionally the Add method
// allows for creating a child subscription. Calling Unsubscribe will close
// the current subscription but will not propagate up to the parent.
// It will traverse recursively all child subcriptions and call Unsubscribe
// on them before settings the subscription state to lifeless.
type Subscriber interface {
	// A Subscriber is also a Subscription.
	Subscription

	// Add will create and return a new child Subscriber with the given
	// callback function. The callback will be called when either the
	// Unsubscribe of the parent or of the returned child subscriber is called.
	// Calling the Unsubscribe method on the child will NOT propagate to the
	// parent! The Unsubscribe will start calling callbacks only after it has
	// set the subscription state to canceled. Even if you call Unsubscribe
	// multiple times, callbacks will only be invoked once.
	Add(func()) Subscriber

	// AddChild will create and return a new child Subscriber. Calling the
	// Unsubscribe method on the child will NOT propagate to the parent!
	// The Unsubscribe will start calling callbacks only after it has
	// set the subscription state to canceled. Even if you call Unsubscribe
	// multiple times, callbacks will only be invoked once.
	AddChild() Subscriber

	// OnUnsubscribe will add the given callback function to the subscriber.
	// The callback will be called when either the Unsubscribe of the parent
	// or of the subscriber itself is called.
	OnUnsubscribe(func())

	// OnWait will register a callback to  call when subscription Wait is called.
	OnWait(func())
}

// New will create and return a new Subscriber.
func New() Subscriber {
	return &subscriber{}
}

const (
	subscribed = iota
	unsubscribed
)

type subscriber struct {
	state int32

	sync.Mutex
	callbacks []func()
	onWait    func()
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

func (s *subscriber) Closed() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscriber) Canceled() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscriber) Wait() {
	s.Lock()
	wait := s.onWait
	s.Unlock()
	if wait != nil {
		wait()
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		s.Add(wg.Done)
		wg.Wait()		
	}
}

func (s *subscriber) Add(callback func()) Subscriber {
	child := &subscriber{callbacks: []func(){callback}}
	s.Lock()
	if atomic.LoadInt32(&s.state) != subscribed {
		child.Unsubscribe()
	} else {
		s.callbacks = append(s.callbacks, child.Unsubscribe)
	}
	s.Unlock()
	return child
}

func (s *subscriber) AddChild() Subscriber {
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
	}
	s.Unlock()
}

func (s *subscriber) OnWait(callback func()) {
	s.Lock()
	s.onWait = callback
	s.Unlock()
}
