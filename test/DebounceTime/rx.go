// Code generated by jig; DO NOT EDIT.

//go:generate jig

package DebounceTime

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:name ConcatInt

// ConcatInt emits the emissions from two or more ObservableInts without interleaving them.
func ConcatInt(observables ...ObservableInt) ObservableInt {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return observables[0].ConcatWith(observables[1:]...)
}

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

//jig:name ObservableInt_ConcatWith

// ConcatWith emits the emissions from two or more ObservableInts without interleaving them.
func (o ObservableInt) ConcatWith(other ...ObservableInt) ObservableInt {
	if len(other) == 0 {
		return o
	}
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			observables	= append([]ObservableInt{}, other...)
			observer	IntObserver
		)
		observer = func(next int, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					var zero int
					observe(zero, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
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

//jig:name Observable_DebounceTime

// DebounceTime only emits the last item of a burst from an Observable if a
// particular timespan has passed without it emitting another item.
func (o Observable) DebounceTime(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var debounce struct {
			sync.Mutex
			runner	scheduler.Runner
			next	interface{}
			done	bool
		}
		debouncer := func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				debounce.Lock()
				debounce.runner = nil
				next := debounce.next
				done := debounce.done
				debounce.Unlock()
				if !done {
					observe(next, nil, false)
				}
			}
		}
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					debounce.Lock()
					debounce.next = next
					if debounce.runner != nil {
						debounce.runner.Cancel()
					}
					debounce.runner = subscribeOn.ScheduleFutureRecursive(duration, debouncer)
					debounce.Unlock()
				} else {
					debounce.Lock()
					debounce.done = true
					debounce.Unlock()
					observe(nil, err, true)
				}
			}
		}
		subscriber.OnUnsubscribe(func() {
			debounce.Lock()
			if debounce.runner != nil {
				debounce.runner.Cancel()
			}
			debounce.Unlock()
		})
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_DebounceTime

// DebounceTime only emits the last item of a burst from an ObservableInt if a
// particular timespan has passed without it emitting another item.
func (o ObservableInt) DebounceTime(duration time.Duration) ObservableInt {
	return o.AsObservable().DebounceTime(duration).AsObservableInt()
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

//jig:name Observable_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o Observable) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next interface{}, err error, done bool) {
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

//jig:name ObservableInt_MapInt

// MapInt transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapInt(project func(int) int) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			var mapped int
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_MergeMapInt

// MergeMapInt transforms the items emitted by an ObservableInt by applying a
// function to each item an returning an ObservableInt. The stream of ObservableInt
// items is then merged into a single stream of Int items using the MergeAll operator.
func (o ObservableInt) MergeMapInt(project func(int) ObservableInt) ObservableInt {
	return o.MapObservableInt(project).MergeAll()
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

//jig:name ObservableInt_MapObservableInt

// MapObservableInt transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapObservableInt(project func(int) ObservableInt) ObservableObservableInt {
	observable := func(observe ObservableIntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			var mapped ObservableInt
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
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
				next.AutoUnsubscribe()(observer, subscribeOn, subscriber)
			} else {
				var zero int
				observer(zero, err, true)
			}
		}
		observers.len += 1
		o.AutoUnsubscribe()(merger, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_AutoUnsubscribeInt

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o ObservableInt) AutoUnsubscribe() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next int, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableObservableInt_AutoUnsubscribeObservableInt

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o ObservableObservableInt) AutoUnsubscribe() ObservableObservableInt {
	observable := func(observe ObservableIntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next ObservableInt, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
