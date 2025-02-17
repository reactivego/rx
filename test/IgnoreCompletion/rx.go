// Code generated by jig; DO NOT EDIT.

//go:generate jig

package IgnoreCompletion

import (
	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

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

//jig:name RangeInt

// RangeInt creates an ObservableInt that emits a range of sequential int values.
// The generated code will do a type conversion from int to int.
func RangeInt(start, count int) ObservableInt {
	end := start + count
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		i := start
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < end {
					observe(int(i), nil, false)
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

//jig:name GoroutineScheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
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

//jig:name NextInt

// NextInt can be called to emit the next value to the IntObserver.
type NextInt func(int)

//jig:name CreateInt

// CreateInt provides a way of creating an ObservableInt from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateInt will be called once
// to implement the observable. It is provided with a NextInt, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func CreateInt(create func(NextInt, Error, Complete, Canceled)) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next int) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				if subscriber.Subscribed() {
					var zero int
					observe(zero, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					var zero int
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

//jig:name Observable_IgnoreCompletion

// IgnoreCompletion only emits items and never completes, neither with Error nor with Complete.
func (o Observable) IgnoreCompletion() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_IgnoreCompletion

// IgnoreCompletion is unknown on the reactivex.io site do we need this?
func (o ObservableInt) IgnoreCompletion() ObservableInt {
	return o.AsObservable().IgnoreCompletion().AsObservableInt()
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

//jig:name Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:name ObservableInt_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Subscribe(observe IntObserver, schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	observer := func(next int, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero int
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
