// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Timeout

import (
	"fmt"
	"sync"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name MakeTrampolineScheduler

func MakeTrampolineScheduler() Scheduler {
	return scheduler.MakeTrampoline()
}

//jig:name Error

// Error signals an error condition.
type Error func(error)

//jig:name Complete

// Complete signals that no more data is to be expected.
type Complete func()

//jig:name NextInt

// NextInt can be called to emit the next value to the IntObserver.
type NextInt func(int)

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

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

//jig:name CreateFutureRecursiveInt

// CreateFutureRecursiveInt provides a way of creating an ObservableInt from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateFutureRecursiveInt will be called
// repeatedly to implement the observable. It is provided with a NextInt, Error
// and Complete function that can be called by the code that implements the
// Observable.
//
// The timeout passed in determines the time before calling the create
// function. The time.Duration returned by the create function determines how
// long CreateFutureRecursiveInt has to wait before calling the create function
// again.
func CreateFutureRecursiveInt(timeout time.Duration, create func(NextInt, Error, Complete) time.Duration) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleFutureRecursive(timeout, func(self func(time.Duration)) {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next int) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				done = true
				if subscriber.Subscribed() {
					var zero int
					observe(zero, err, true)
				}
			}
			c := func() {
				done = true
				if subscriber.Subscribed() {
					var zero int
					observe(zero, nil, true)
				}
			}
			timeout = create(n, e, c)
			if !done && subscriber.Subscribed() {
				self(timeout)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name DeferInt

// DeferInt does not create the ObservableInt until the observer subscribes.
// It creates a fresh ObservableInt for each subscribing observer. Use it to
// create observables that maintain separate state per subscription.
func DeferInt(factory func() ObservableInt) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name NeverInt

// NeverInt creates an ObservableInt that emits no items and does't terminate.
func NeverInt() ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
	}
	return observable
}

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

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name Observable_Serialize

// Serialize forces an Observable to make serialized calls and to be
// well-behaved.
func (o Observable) Serialize() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var observer struct {
			sync.Mutex
			done	bool
		}
		serializer := func(next interface{}, err error, done bool) {
			observer.Lock()
			defer observer.Unlock()
			if !observer.done {
				observer.done = done
				observe(next, err, done)
			}
		}
		o(serializer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_Timeout

// TimeoutOccured is delivered to an observer if the stream times out.
const TimeoutOccured = RxError("timeout occured")

// Timeout mirrors the source Observable, but issues an error notification if a
// particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o Observable) Timeout(due time.Duration) Observable {
	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var timeout struct {
			sync.Mutex
			at		time.Time
			occurred	bool
		}
		timeout.at = subscribeOn.Now().Add(due)
		timer := subscribeOn.ScheduleFutureRecursive(due, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					due := timeout.at.Sub(subscribeOn.Now())
					if due > 0 {
						self(due)
					} else {
						timeout.occurred = true
						timeout.Unlock()
						observe(nil, TimeoutOccured, true)
						timeout.Lock()
					}
				}
				timeout.Unlock()
			}
		})
		subscriber.OnUnsubscribe(timer.Cancel)
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					now := subscribeOn.Now()
					if now.Before(timeout.at) {
						timeout.at = now.Add(due)
						timeout.occurred = done
						observe(next, err, done)
					}
				}
				timeout.Unlock()
			}
		}
		o(observer, subscribeOn, subscriber)
	})
	return observable
}

//jig:name ObservableInt_Timeout

// Timeout mirrors the source ObservableInt, but issues an error notification if
// a particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o ObservableInt) Timeout(timeout time.Duration) ObservableInt {
	return o.AsObservable().Timeout(timeout).AsObservableInt()
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

//jig:name ObservableInt_SubscribeOn

// SubscribeOn specifies the scheduler an ObservableInt should use when it is
// subscribed to.
func (o ObservableInt) SubscribeOn(scheduler Scheduler) ObservableInt {
	observable := func(observe IntObserver, _ Scheduler, subscriber Subscriber) {
		if scheduler.IsConcurrent() {
			subscriber.OnWait(nil)
		} else {
			subscriber.OnWait(scheduler.Wait)
		}
		o(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name ObservableInt_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
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
