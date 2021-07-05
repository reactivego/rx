// Code generated by jig; DO NOT EDIT.

//go:generate jig

package MergeAll

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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

//jig:name NextObservableString

// NextObservableString can be called to emit the next value to the IntObserver.
type NextObservableString func(ObservableString)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

//jig:name ObservableStringObserver

// ObservableStringObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableStringObserver func(next ObservableString, err error, done bool)

//jig:name ObservableObservableString

// ObservableObservableString is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservableString func(ObservableStringObserver, Scheduler, Subscriber)

//jig:name CreateObservableString

// CreateObservableString provides a way of creating an ObservableObservableString from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateObservableString will be called once
// to implement the observable. It is provided with a NextObservableString, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func CreateObservableString(create func(NextObservableString, Error, Complete, Canceled)) ObservableObservableString {
	observable := func(observe ObservableStringObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next ObservableString) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				if subscriber.Subscribed() {
					var zero ObservableString
					observe(zero, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					var zero ObservableString
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

//jig:name JustString

// JustString creates an ObservableString that emits a particular item.
func JustString(element string) ObservableString {
	observable := func(observe StringObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(element, nil, false)
			}
			if subscriber.Subscribed() {
				var zero string
				observe(zero, nil, true)
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

//jig:name ObservableIntObserver

// ObservableIntObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableIntObserver func(next ObservableInt, err error, done bool)

//jig:name ObservableObservableInt

// ObservableObservableInt is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservableInt func(ObservableIntObserver, Scheduler, Subscriber)

//jig:name FromObservableInt

// FromObservableInt creates an ObservableObservableInt from multiple ObservableInt values passed in.
func FromObservableInt(slice ...ObservableInt) ObservableObservableInt {
	observable := func(observe ObservableIntObserver, scheduler Scheduler, subscriber Subscriber) {
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
					var zero ObservableInt
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableObservableString_MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableString) MergeAll() ObservableString {
	observable := func(observe StringObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done	bool
			len	int32
		}
		observer := func(next string, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					var zero string
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zero string
						observe(zero, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableString, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				var zero string
				observer(zero, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1
				o(merger, subscribeOn, subscriber)
			}
		})
	}
	return observable
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

//jig:name ObservableInt_Take

// Take emits only the first n items emitted by an ObservableInt.
func (o ObservableInt) Take(n int) ObservableInt {
	return o.AsObservable().Take(n).AsObservableInt()
}

//jig:name ObservableObservableInt_MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableInt) MergeAll() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done	bool
			len	int32
		}
		observer := func(next int, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					var zero int
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zero int
						observe(zero, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableInt, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				var zero int
				observer(zero, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1
				o(merger, subscribeOn, subscriber)
			}
		})
	}
	return observable
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

//jig:name ObservableString_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableString) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
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

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

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
