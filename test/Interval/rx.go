// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Interval

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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

//jig:name Float32Observer

// Float32Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Float32Observer func(next float32, err error, done bool)

//jig:name ObservableFloat32

// ObservableFloat32 is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableFloat32 func(Float32Observer, Scheduler, Subscriber)

//jig:name IntervalFloat32

// IntervalFloat32 creates an ObservableFloat32 that emits a sequence of integers spaced
// by a particular time interval. First integer is not emitted immediately, but
// only after the first time interval has passed. The generated code will do a type
// conversion from int to float32.
func IntervalFloat32(interval time.Duration) ObservableFloat32 {
	observable := func(observe Float32Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(float32(i), nil, false)
				i++
				if subscriber.Subscribed() {
					self(interval)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name GoroutineScheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
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

//jig:name Interval

// Interval creates an Observable that emits a sequence of integers spaced
// by a particular time interval. First integer is not emitted immediately, but
// only after the first time interval has passed. The generated code will do a type
// conversion from int to interface{}.
func Interval(interval time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(interface{}(i), nil, false)
				i++
				if subscriber.Subscribed() {
					self(interval)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name IntObserver

// IntObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type IntObserver func(next int, err error, done bool)

//jig:name ObservableInt

// ObservableInt is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableInt func(IntObserver, Scheduler, Subscriber)

//jig:name IntervalInt

// IntervalInt creates an ObservableInt that emits a sequence of integers spaced
// by a particular time interval. First integer is not emitted immediately, but
// only after the first time interval has passed. The generated code will do a type
// conversion from int to int.
func IntervalInt(interval time.Duration) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(int(i), nil, false)
				i++
				if subscriber.Subscribed() {
					self(interval)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

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

//jig:name Observable_Take

// Take emits only the first n items emitted by an Observable.
func (o Observable) Take(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		taken := 0
		observer := func(next interface{}, err error, done bool) {
			if taken < n {
				observe(next, err, done)
				if !done {
					taken++
					if taken >= n {
						observe(nil, nil, true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableFloat32_Take

// Take emits only the first n items emitted by an ObservableFloat32.
func (o ObservableFloat32) Take(n int) ObservableFloat32 {
	return o.AsObservable().Take(n).AsObservableFloat32()
}

//jig:name Observable_SubscribeOn

// SubscribeOn specifies the scheduler an Observable should use when it is
// subscribed to.
func (o Observable) SubscribeOn(scheduler Scheduler) Observable {
	observable := func(observe Observer, _ Scheduler, subscriber Subscriber) {
		if scheduler.IsConcurrent() {
			subscriber.OnWait(nil)
		} else {
			subscriber.OnWait(scheduler.Wait)
		}
		o(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name ObservableInt_Take

// Take emits only the first n items emitted by an ObservableInt.
func (o ObservableInt) Take(n int) ObservableInt {
	return o.AsObservable().Take(n).AsObservableInt()
}

//jig:name ObservableFloat32_Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a serial scheduler created with NewScheduler().
func (o ObservableFloat32) Wait() error {
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next float32, err error, done bool) {
		if done {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:name Observable_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a serial scheduler created with NewScheduler().
func (o Observable) Subscribe(observe Observer, schedulers ...Scheduler) Subscription {
	subscriber := NewSubscriber()
	schedulers = append(schedulers, NewScheduler())
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero interface{}
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

//jig:name ObservableInt_MapString

// MapString transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapString(project func(int) string) ObservableString {
	observable := func(observe StringObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			var mapped string
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableFloat32_AsObservable

// AsObservable turns a typed ObservableFloat32 into an Observable of interface{}.
func (o ObservableFloat32) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next float32, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_AsObservable

// AsObservable turns a typed ObservableInt into an Observable of interface{}.
func (o ObservableInt) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name StringObserver

// StringObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type StringObserver func(next string, err error, done bool)

//jig:name ObservableString

// ObservableString is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableString func(StringObserver, Scheduler, Subscriber)

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name TypecastFailed

// ErrTypecast is delivered to an observer if the generic value cannot be
// typecast to a specific type.
const TypecastFailed = RxError("typecast failed")

//jig:name Observable_AsObservableFloat32

// AsObservableFloat32 turns an Observable of interface{} into an ObservableFloat32.
// If during observing a typecast fails, the error ErrTypecastToFloat32 will be
// emitted.
func (o Observable) AsObservableFloat32() ObservableFloat32 {
	observable := func(observe Float32Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextFloat32, ok := next.(float32); ok {
					observe(nextFloat32, err, done)
				} else {
					var zero float32
					observe(zero, TypecastFailed, true)
				}
			} else {
				var zero float32
				observe(zero, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_AsObservableInt

// AsObservableInt turns an Observable of interface{} into an ObservableInt.
// If during observing a typecast fails, the error ErrTypecastToInt will be
// emitted.
func (o Observable) AsObservableInt() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextInt, ok := next.(int); ok {
					observe(nextInt, err, done)
				} else {
					var zero int
					observe(zero, TypecastFailed, true)
				}
			} else {
				var zero int
				observe(zero, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
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

//jig:name NewScheduler

func NewScheduler() Scheduler {
	return scheduler.New()
}

//jig:name ObservableString_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a serial scheduler created with NewScheduler().
func (o ObservableString) Println(a ...interface{}) error {
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next string, err error, done bool) {
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
