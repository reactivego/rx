// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Last

import (
	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

// NewSubscriber creates a new subscriber.
func NewSubscriber() Subscriber {
	return subscriber.New()
}

//jig:name IntObserveFunc

// IntObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type IntObserveFunc func(next int, err error, done bool)

//jig:name zeroInt

var zeroInt int

//jig:name ObservableInt

// ObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableInt func(IntObserveFunc, Scheduler, Subscriber)

//jig:name FromSliceInt

// FromSliceInt creates an ObservableInt from a slice of int values passed in.
func FromSliceInt(slice []int) ObservableInt {
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
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
					observe(zeroInt, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name FromInts

// FromInts creates an ObservableInt from multiple int values passed in.
func FromInts(slice ...int) ObservableInt {
	return FromSliceInt(slice)
}

//jig:name ObservableLast

// Last emits only the last item emitted by an Observable.
func (o Observable) Last() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		have := false
		var last interface{}
		observer := func(next interface{}, err error, done bool) {
			if done {
				if have {
					observe(last, nil, false)
				}
				observe(nil, err, true)
			} else {
				last = next
				have = true
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntLast

// Last emits only the last item emitted by an ObservableInt.
func (o ObservableInt) Last() ObservableInt {
	return o.AsObservable().Last().AsObservableInt()
}

//jig:name ObserveFunc

// ObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type ObserveFunc func(next interface{}, err error, done bool)

//jig:name zero

var zero interface{}

//jig:name Observable

// Observable is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type Observable func(ObserveFunc, Scheduler, Subscriber)

//jig:name ObservableIntAsObservable

// AsObservable turns a typed ObservableInt into an Observable of interface{}.
func (o ObservableInt) AsObservable() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Schedulers

func TrampolineScheduler() Scheduler	{ return scheduler.Trampoline }

func GoroutineScheduler() Scheduler	{ return scheduler.Goroutine }

//jig:name ObservableIntToSingle

// ToSingle blocks until the ObservableInt emits exactly one value or an error.
// The value and any error are returned.
func (o ObservableInt) ToSingle() (entry int, err error) {
	o = o.Single()
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next int, e error, done bool) {
		if !done {
			entry = next
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name ErrTypecastToInt

// ErrTypecastToInt is delivered to an observer if the generic value cannot be
// typecast to int.
const ErrTypecastToInt = RxError("typecast to int failed")

//jig:name ObservableAsObservableInt

// AsInt turns an Observable of interface{} into an ObservableInt. If during
// observing a typecast fails, the error ErrTypecastToInt will be emitted.
func (o Observable) AsObservableInt() ObservableInt {
	observable := func(observe IntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextInt, ok := next.(int); ok {
					observe(nextInt, err, done)
				} else {
					observe(zeroInt, ErrTypecastToInt, true)
				}
			} else {
				observe(zeroInt, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntSingle

// Single enforces that the observableInt sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o ObservableInt) Single() ObservableInt {
	return o.AsObservable().Single().AsObservableInt()
}

//jig:name ObservableSingle

// Single enforces that the observable sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o Observable) Single() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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
							observe(nil, RxError("expected one value, got none"), true)
						}
					}
				} else {
					count++
					if count == 1 {
						latest = next
					} else {
						observe(nil, RxError("expected one value, got multiple"), true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
