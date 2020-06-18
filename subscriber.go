package subscriber

import (
	"sync"

	"sync/atomic"
)

// Error type to allow constant error declarations.
type Error string

// Error method implements Error interface.
func (e Error) Error() string {
	return string(e)
}

// Unsubscribed is the error returned by wait when the Unsubscribe method
// is called on the subscription.
const ErrUnsubscribed = Error("subscriber unsubscribed")

// Subscription is an interface that allows code to monitor and control a
// subscription it received.
type Subscription interface {
	// Subscribed returns true when the subscription is currently active.
	Subscribed() bool

	// Unsubscribe will cancel the subscription when it is still active.
	// Subsequently, it will call Unsubscribe on all child subscriptions added
	// via Add along with all methods added through OnUnsubscribe. Unsubscribe
	// can be safely called on a subscription that is no longer active and
	// performs no operation in that case. There is no need to check if the
	// subscription is still active before calling Unsubscribe. When the
	// subscription is canceled by calling Unsubscribe a call to the Wait
	// method will return the error ErrUnsubscribed.
	Unsubscribe()

	// Canceled returns true when the subscription state is canceled.
	Canceled() bool

	// Wait will by default block the calling goroutine and wait for the
	// Unsubscribe method to be called on this subscription.
	// However, when OnWait was called with a callback wait function it will
	// call that instead. Calling Wait on a subscription that has already been
	// canceled will return immediately. If the subscriber was canceled by
	// calling Unsubscribe, then the error returned is ErrUnsubscribed.
	// If the subscriber was terminated by calling Done, then the error
	// returned here is the one passed to Done.
	Wait() error
}

// Subscriber is a Subscription with management functionality.
type Subscriber interface {
	// A Subscriber is also a Subscription.
	Subscription

	// Add will create and return a new child Subscriber with the given
	// callback functions. The callbacks will be called when either the
	// Unsubscribe of the parent or of the returned child subscriber is called.
	// Calling the Unsubscribe method on the child will NOT propagate to the
	// parent! The Unsubscribe will start calling callbacks only after it has
	// set the subscription state to canceled. Even if you call Unsubscribe
	// multiple times, callbacks will only be invoked once.
	Add(callbacks ...func()) Subscriber

	// OnUnsubscribe will add the given callback function to the subscriber.
	// The callback will be called when either the Unsubscribe of the parent
	// or of the subscriber itself is called. If the subscription was already
	// canceled, then the callback function will just be called immediately.
	OnUnsubscribe(callback func())

	// OnWait will register a callback to  call when subscription Wait is called.
	OnWait(callback func())

	// Done will set the error internally and then cancel the subscription by
	// calling the Unsubscribe method. A nil value for error indicates success.
	Done(err error)

	// Error returns the error set by calling the Done(err) method. As long as
	// the subscriber is still subscribed Error will return nil.
	Error() error
}

// New will create and return a new Subscriber.
func New(callbacks ...func()) Subscriber {
	return &subscriber{callbacks: callbacks, err: ErrUnsubscribed}
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
	err       error
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

func (s *subscriber) Wait() error {
	s.Lock()
	wait := s.onWait
	s.Unlock()
	if wait != nil {
		wait()
	}
	if atomic.LoadInt32(&s.state) == subscribed {
		var wg sync.WaitGroup
		wg.Add(1)
		s.OnUnsubscribe(wg.Done)
		wg.Wait()
	}
	return s.Error()
}

func (s *subscriber) Add(callbacks ...func()) Subscriber {
	child := New(callbacks...)
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

func (s *subscriber) OnWait(callback func()) {
	s.Lock()
	s.onWait = callback
	s.Unlock()
}

func (s *subscriber) Done(err error) {
	s.Lock()
	s.err = err
	s.Unlock()
	s.Unsubscribe()
}

func (s *subscriber) Error() error {
	s.Lock()
	err := s.err
	s.Unlock()
	if atomic.LoadInt32(&s.state) == subscribed {
		err = nil
	}
	return err
}
