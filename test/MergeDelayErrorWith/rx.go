// Code generated by jig; DO NOT EDIT.

//go:generate jig

package MergeDelayErrorWith

import (
	"sync"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

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

//jig:name ObservableInt_MergeDelayErrorWith

// MergeDelayError combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func (o ObservableInt) MergeDelayErrorWith(other ...ObservableInt) ObservableInt {
	if len(other) == 0 {
		return o
	}
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			len	int
			err	error
		}
		observer := func(next int, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !done {
				observe(next, nil, false)
			} else {
				if err != nil {
					observers.err = err
				}
				if observers.len--; observers.len == 0 {
					var zero int
					observe(zero, observers.err, true)
				}
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1 + len(other)
				o(observer, subscribeOn, subscriber)
				for _, o := range other {
					if !subscriber.Subscribed() {
						return
					}
					o(observer, subscribeOn, subscriber)
				}
			}
		})
	}
	return observable
}

//jig:name ObservableInt_ToSlice

// ToSlice collects all values from the ObservableInt into an slice. The
// complete slice and any error are returned.
// ToSlice uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) ToSlice() (slice []int, err error) {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
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
