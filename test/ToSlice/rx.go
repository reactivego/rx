// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ToSlice

import (
	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber = subscriber.Subscriber

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

//jig:name Range

// Range creates an ObservableInt that emits a range of sequential integers.
func Range(start, count int) ObservableInt {
	end := start + count
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := start
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < end {
					observe(i, nil, false)
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

//jig:name Throw

// Throw creates an Observable that emits no items and terminates with an
// error.
func Throw(err error) Observable {
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zero, err, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name Schedulers

func TrampolineScheduler() Scheduler {
	return scheduler.Trampoline
}

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

//jig:name ObservableIntToSlice

// ToSlice collects all values from the ObservableInt into an slice. The
// complete slice and any error are returned.
func (o ObservableInt) ToSlice() (slice []int, err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next int, e error, done bool) {
		if !done {
			slice = append(slice, next)
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

//jig:name ObservableToSlice

// ToSlice collects all values from the Observable into an slice. The
// complete slice and any error are returned.
func (o Observable) ToSlice() (slice []interface{}, err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next interface{}, e error, done bool) {
		if !done {
			slice = append(slice, next)
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
