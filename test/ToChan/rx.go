// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ToChan

import (
	"sync"
	"sync/atomic"

	"github.com/reactivego/scheduler"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is a Subscription with management functionality.
type Subscriber interface {
	// A Subscriber is also a Subscription.
	Subscription

	// Add will create and return a new child Subscriber setup in such a way that
	// calling Unsubscribe on the parent will also call Unsubscribe on the child.
	// Calling the Unsubscribe method on the child will NOT propagate to the
	// parent!
	Add() Subscriber

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

//jig:name Observer

// Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Observer func(next interface{}, err error, done bool)

//jig:name Observable

// Observable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type Observable func(Observer, Scheduler, Subscriber)

//jig:name From

// From creates an Observable from multiple interface{} values passed in.
func From(slice ...interface{}) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					var zero interface{}
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name NewSubscriber

// New will create and return a new Subscriber.
func NewSubscriber() Subscriber {
	return &subscriber{err: ErrUnsubscribed}
}

// Unsubscribed is the error returned by wait when the Unsubscribe method
// is called on the subscription.
const ErrUnsubscribed = RxError("subscriber unsubscribed")

const (
	subscribed	= iota
	unsubscribed
)

type subscriber struct {
	state	int32

	sync.Mutex
	callbacks	[]func()
	onWait		func()
	err		error
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

func (s *subscriber) Add() Subscriber {
	child := NewSubscriber()
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

//jig:name Range

// Range creates an Observable that emits a range of sequential int values.
// The generated code will do a type conversion from int to interface{}.
func Range(start, count int) Observable {
	end := start + count
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		i := start
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < end {
					observe(interface{}(i), nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					var zero interface{}
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Error

// Error signals an error condition.
type Error func(error)

//jig:name Complete

// Complete signals that no more data is to be expected.
type Complete func()

//jig:name Canceled

// Canceled returns true when the observer has unsubscribed.
type Canceled func() bool

//jig:name Next

// Next can be called to emit the next value to the IntObserver.
type Next func(interface{})

//jig:name Create

// Create provides a way of creating an Observable from
// scratch by calling observer methods programmatically.
//
// The create function provided to Create will be called once
// to implement the observable. It is provided with a Next, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func Create(create func(Next, Error, Complete, Canceled)) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next interface{}) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				if subscriber.Subscribed() {
					var zero interface{}
					observe(zero, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					var zero interface{}
					observe(zero, nil, true)
				}
			}
			x := func() bool {
				return !subscriber.Subscribed()
			}
			create(n, e, c, x)
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name Subscription

// Subscription is an interface that allows code to monitor and control a
// subscription it received.
type Subscription interface {
	// Subscribed returns true when the subscription is currently active.
	Subscribed() bool

	// Unsubscribe will do nothing if the subscription is not active. If the
	// state is still active however, it will be changed to canceled.
	// Subsequently, it will call Unsubscribe on all child subscriptions added
	// through Add, along with all methods added through OnUnsubscribe. When the
	// subscription is canceled by calling Unsubscribe a call to the Wait method
	// will return the error ErrUnsubscribed.
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

//jig:name Observable_ToChan

// ToChan returns a channel that emits interface{} values. If the source
// observable does not emit values but emits an error or complete, then the
// returned channel will enit any error and then close without emitting any
// values.
//
// ToChan uses the public scheduler.Goroutine variable for scheduling, because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To be able to cancel ToChan, you will need to
// create a subscriber yourself and pass it to ToChan as an argument.
func (o Observable) ToChan(subscribers ...Subscriber) <-chan interface{} {
	subscribers = append(subscribers, NewSubscriber())
	scheduler := scheduler.Goroutine
	donech := make(chan struct{})
	nextch := make(chan interface{})
	const (
		idle	= iota
		busy
		closed
	)
	state := int32(idle)
	observer := func(next interface{}, err error, done bool) {
		if atomic.CompareAndSwapInt32(&state, idle, busy) {
			if err != nil {
				next = err
			}
			if !done || err != nil {
				select {
				case <-donech:
					atomic.StoreInt32(&state, closed)
				default:
					select {
					case <-donech:
						atomic.StoreInt32(&state, closed)
					case nextch <- next:
					}
				}
			}
			if done {
				atomic.StoreInt32(&state, closed)
				subscribers[0].Done(err)
			}
			if !atomic.CompareAndSwapInt32(&state, busy, idle) {
				close(nextch)
			}
		}
	}
	subscribers[0].OnUnsubscribe(func() {
		close(donech)
		if atomic.CompareAndSwapInt32(&state, busy, closed) {
			return
		}
		if atomic.CompareAndSwapInt32(&state, idle, closed) {
			close(nextch)
			return
		}
	})
	o(observer, scheduler, subscribers[0])
	return nextch
}
