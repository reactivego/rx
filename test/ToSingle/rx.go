// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ToSingle

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

//jig:name FromInt

// FromInt creates an ObservableInt from multiple int values passed in.
func FromInt(slice ...int) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
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
					var zero int
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name EmptyInt

// EmptyInt creates an Observable that emits no items but terminates normally.
func EmptyInt() ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				var zero int
				observe(zero, nil, true)
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

//jig:name ObservableInt_ToSingle

// ToSingle blocks until the ObservableInt emits exactly one value or an error.
// The value and any error are returned.
// ToSingle uses a serial scheduler created with NewScheduler().
func (o ObservableInt) ToSingle() (entry int, err error) {
	o = o.Single()
	subscriber := NewSubscriber()
	scheduler := NewScheduler()
	observer := func(next int, err error, done bool) {
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

//jig:name ObservableInt_Single

// Single enforces that the observableInt sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o ObservableInt) Single() ObservableInt {
	return o.AsObservable().Single().AsObservableInt()
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

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

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

//jig:name ErrSingle

const DidNotEmitValue = RxError("expected one value, got none")

const EmittedMultipleValues = RxError("expected one value, got multiple")

//jig:name Observable_Single

// Single enforces that the observable sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing, this is reported as an error to the observer.
func (o Observable) Single() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			count	int
			latest	interface{}
		)
		observer := func(next interface{}, err error, done bool) {
			if count < 2 {
				if done {
					if err != nil {
						observe(nil, err, true)
					} else {
						if count == 1 {
							observe(latest, nil, false)
							observe(nil, nil, true)
						} else {
							observe(nil, DidNotEmitValue, true)
						}
					}
				} else {
					count++
					if count == 1 {
						latest = next
					} else {
						observe(nil, EmittedMultipleValues, true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name TypecastFailed

// ErrTypecast is delivered to an observer if the generic value cannot be
// typecast to a specific type.
const TypecastFailed = RxError("typecast failed")

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
