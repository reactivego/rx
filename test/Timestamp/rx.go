// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Timestamp

import (
	"fmt"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name TimestampTime

type TimestampTime struct {
	Value		Time
	Timestamp	time.Time
}

//jig:name Time

type Time = time.Time

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

//jig:name TimeObserver

// TimeObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type TimeObserver func(next Time, err error, done bool)

//jig:name ObservableTime

// ObservableTime is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableTime func(TimeObserver, Scheduler, Subscriber)

//jig:name Ticker

// Ticker creates an ObservableTime that emits a sequence of timestamps after
// an initialDelay has passed. Subsequent timestamps are emitted using a
// schedule of intervals passed in. If only the initialDelay is given, Ticker
// will emit only once.
func Ticker(initialDelay time.Duration, intervals ...time.Duration) ObservableTime {
	observable := func(observe TimeObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(initialDelay, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				if i == 0 || (i > 0 && len(intervals) > 0) {
					observe(subscribeOn.Now(), nil, false)
				}
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						self(intervals[i%len(intervals)])
					} else {
						if i == 0 {
							self(0)
						} else {
							var zero time.Time
							observe(zero, nil, true)
						}
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableTime_Timestamp

// Timestamp attaches a timestamp to each item emitted by an observable
// indicating when it was emitted.
func (o ObservableTime) Timestamp() ObservableTimestampTime {
	observable := func(observe TimestampTimeObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next Time, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					observe(TimestampTime{next, subscribeOn.Now()}, nil, false)
				} else {
					var zero TimestampTime
					observe(zero, err, done)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name TimestampTimeObserver

// TimestampTimeObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type TimestampTimeObserver func(next TimestampTime, err error, done bool)

//jig:name ObservableTimestampTime

// ObservableTimestampTime is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableTimestampTime func(TimestampTimeObserver, Scheduler, Subscriber)

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

//jig:name ObservableTimestampTime_Take

// Take emits only the first n items emitted by an ObservableTimestampTime.
func (o ObservableTimestampTime) Take(n int) ObservableTimestampTime {
	return o.AsObservable().Take(n).AsObservableTimestampTime()
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

//jig:name ObservableTimestampTime_AsObservable

// AsObservable turns a typed ObservableTimestampTime into an Observable of interface{}.
func (o ObservableTimestampTime) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next TimestampTime, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableTimestampTime_All

// All determines whether all items emitted by an ObservableTimestampTime meet some
// criteria.
//
// Pass a predicate function to the All operator that accepts an item emitted
// by the source ObservableTimestampTime and returns a boolean value based on an
// evaluation of that item. All returns an ObservableBool that emits a single
// boolean value: true if and only if the source ObservableTimestampTime terminates
// normally and every item emitted by the source ObservableTimestampTime evaluated as
// true according to this predicate; false if any item emitted by the source
// ObservableTimestampTime evaluates as false according to this predicate.
func (o ObservableTimestampTime) All(predicate func(next TimestampTime) bool) ObservableBool {
	condition := func(next interface{}) bool {
		return predicate(next.(TimestampTime))
	}
	return o.AsObservable().All(condition)
}

//jig:name BoolObserver

// BoolObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type BoolObserver func(next bool, err error, done bool)

//jig:name ObservableBool

// ObservableBool is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableBool func(BoolObserver, Scheduler, Subscriber)

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name TypecastFailed

// ErrTypecast is delivered to an observer if the generic value cannot be
// typecast to a specific type.
const TypecastFailed = RxError("typecast failed")

//jig:name Observable_AsObservableTimestampTime

// AsObservableTimestampTime turns an Observable of interface{} into an ObservableTimestampTime.
// If during observing a typecast fails, the error ErrTypecastToTimestampTime will be
// emitted.
func (o Observable) AsObservableTimestampTime() ObservableTimestampTime {
	observable := func(observe TimestampTimeObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextTimestampTime, ok := next.(TimestampTime); ok {
					observe(nextTimestampTime, err, done)
				} else {
					var zero TimestampTime
					observe(zero, TypecastFailed, true)
				}
			} else {
				var zero TimestampTime
				observe(zero, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_All

// All determines whether all items emitted by an Observable meet some
// criteria.
//
// Pass a predicate function to the All operator that accepts an item emitted
// by the source Observable and returns a boolean value based on an
// evaluation of that item. All returns an ObservableBool that emits a single
// boolean value: true if and only if the source Observable terminates
// normally and every item emitted by the source Observable evaluated as
// true according to this predicate; false if any item emitted by the source
// Observable evaluates as false according to this predicate.
func (o Observable) All(predicate func(next interface{}) bool) ObservableBool {
	observable := func(observe BoolObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			switch {
			case !done:
				if !predicate(next) {
					observe(false, nil, false)
					observe(false, nil, true)
				}
			case err != nil:
				observe(false, err, true)
			default:
				observe(true, nil, false)
				observe(false, nil, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableBool_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableBool) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next bool, err error, done bool) {
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
