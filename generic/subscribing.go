package rx

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/reactivego/scheduler"
)

//jig:template Subscription

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

//jig:template Subscriber

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

//jig:template NewSubscriber
//jig:needs Subscriber

// New will create and return a new Subscriber.
func NewSubscriber() Subscriber {
	return &subscriber{err: ErrUnsubscribed}
}

// Unsubscribed is the error returned by wait when the Unsubscribe method
// is called on the subscription.
const ErrUnsubscribed = RxError("subscriber unsubscribed")

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

//jig:template Observable<Foo> AutoUnsubscribe<Foo>

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o ObservableFoo) AutoUnsubscribe() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next foo, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Println
//jig:needs Scheduler, Subscriber

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a serial scheduler created with NewScheduler().
func (o ObservableFoo) Println(a ...interface{}) error {
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next foo, err error, done bool) {
		if !done {
			fmt.Println(append(a, next)...)
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:template Observable<Foo> Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a serial scheduler created with NewScheduler().
func (o ObservableFoo) Wait() error {
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next foo, err error, done bool) {
		if done {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:template Observable<Foo> ToSlice

// ToSlice collects all values from the ObservableFoo into an slice. The
// complete slice and any error are returned.
// ToSlice uses a serial scheduler created with NewScheduler().
func (o ObservableFoo) ToSlice() (slice []foo, err error) {
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next foo, err error, done bool) {
		if !done {
			slice = append(slice, next)
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	err = subscriber.Wait()
	return
}

//jig:template Observable<Foo> ToSingle

// ToSingle blocks until the ObservableFoo emits exactly one value or an error.
// The value and any error are returned.
// ToSingle uses a serial scheduler created with NewScheduler().
func (o ObservableFoo) ToSingle() (entry foo, err error) {
	o = o.Single()
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next foo, err error, done bool) {
		if !done {
			entry = next
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	err = subscriber.Wait()
	return
}

//jig:template Observable ToChan

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
		idle = iota
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

//jig:template Observable<Foo> ToChan
//jig:required-vars Foo

// ToChan returns a channel that emits foo values. If the source observable does
// not emit values but emits an error or complete, then the returned channel
// will close without emitting any values.
//
// ToChan uses the public scheduler.Goroutine variable for scheduling, because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To be able to cancel ToChan, you will need to
// create a subscriber yourself and pass it to ToChan as an argument.
func (o ObservableFoo) ToChan(subscribers ...Subscriber) <-chan foo {
	subscribers = append(subscribers, NewSubscriber())
	scheduler := scheduler.Goroutine
	donech := make(chan struct{})
	nextch := make(chan foo)
	const (
		idle = iota
		busy
		closed
	)
	state := int32(idle)
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&state, idle, busy) {
			if !done {
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
			} else {
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

//jig:template Observable<Foo> Subscribe
//jig:needs Subscription

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a serial scheduler created with NewScheduler().
func (o ObservableFoo) Subscribe(observe FooObserver, schedulers ...Scheduler) Subscription {
	subscriber := NewSubscriber()
	schedulers = append(schedulers, NewScheduler())
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero foo
			observe(zero, err, true)
			subscriber.Done(err)
		}
	}
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	o(observer, schedulers[0], subscriber)
	return subscriber
}

//jig:template Println<Foo>
//jig:needs <Foo>Observer

func PrintlnFoo(a ...interface{}) FooObserver {
	observer := func(next foo, err error, done bool) {
		if !done {
			fmt.Println(append(a, next)...)
		}
	}
	return observer
}
