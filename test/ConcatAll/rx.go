// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ConcatAll

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

//jig:name JustInt

// JustInt creates an ObservableInt that emits a particular item.
func JustInt(element int) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(element, nil, false)
			}
			if subscriber.Subscribed() {
				var zero int
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
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

//jig:name Observable_Delay

// Delay shifts an emission from an Observable forward in time by a particular
// amount of time. The relative time intervals between emissions are preserved.
func (o Observable) Delay(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		type emission struct {
			at	time.Time
			next	interface{}
			err	error
			done	bool
		}
		var delay struct {
			sync.Mutex
			emissions	[]emission
		}
		delayer := subscribeOn.ScheduleFutureRecursive(duration, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				delay.Lock()
				for _, entry := range delay.emissions {
					delay.Unlock()
					due := entry.at.Sub(subscribeOn.Now())
					if due > 0 {
						self(due)
						return
					}
					observe(entry.next, entry.err, entry.done)
					if entry.done || !subscriber.Subscribed() {
						return
					}
					delay.Lock()
					delay.emissions = delay.emissions[1:]
				}
				delay.Unlock()
				self(duration)
			}
		})
		subscriber.OnUnsubscribe(delayer.Cancel)
		observer := func(next interface{}, err error, done bool) {
			delay.Lock()
			delay.emissions = append(delay.emissions, emission{subscribeOn.Now().Add(duration), next, err, done})
			delay.Unlock()
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_Delay

// Delay shifts an emission from an Observable forward in time by a particular
// amount of time. The relative time intervals between emissions are preserved.
func (o ObservableInt) Delay(duration time.Duration) ObservableInt {
	return o.AsObservable().Delay(duration).AsObservableInt()
}

//jig:name ObservableInt_DoOnComplete

// DoOnComplete calls a function when the stream completes.
func (o ObservableInt) DoOnComplete(f func()) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			if err == nil && done {
				f()
			}
			observe(next, err, done)
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

//jig:name ObservableObservableInt_ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableInt) ConcatAll() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var concat struct {
			sync.Mutex
			observables	[]ObservableInt
			observer	IntObserver
			subscriber	Subscriber
		}

		var source struct {
			observer	ObservableIntObserver
			subscriber	Subscriber
		}

		concat.observer = func(next int, err error, done bool) {
			concat.Lock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(concat.observables) == 0 {
					if !source.subscriber.Subscribed() {
						var zero int
						observe(zero, nil, true)
					}
					concat.observables = nil
				} else {
					observable := concat.observables[0]
					concat.observables = concat.observables[1:]
					observable(concat.observer, subscribeOn, subscriber)
				}
			}
			concat.Unlock()
		}

		source.observer = func(next ObservableInt, err error, done bool) {
			if !done {
				concat.Lock()
				initial := concat.observables == nil
				concat.observables = append(concat.observables, next)
				concat.Unlock()
				if initial {
					var zero int
					concat.observer(zero, nil, true)
				}
			} else {
				concat.Lock()
				initial := concat.observables == nil
				source.subscriber.Done(err)
				concat.Unlock()
				if initial || err != nil {
					var zero int
					concat.observer(zero, err, true)
				}
			}
		}
		source.subscriber = subscriber.Add()
		o(source.observer, subscribeOn, source.subscriber)
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

//jig:name ObservableInt_Do

// Do calls a function for each next value passing through the observable.
func (o ObservableInt) Do(f func(next int)) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			if !done {
				f(next)
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
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

//jig:name ObservableInt_Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Wait() error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
		if done {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}
