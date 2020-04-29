// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Make

import (
	"fmt"

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

//jig:name MakeIntFunc

// MakeIntFunc is the signature of a function that can be passed to MakeInt
// to implement an ObservableInt.
type MakeIntFunc func(Next func(int), Error func(error), Complete func())

//jig:name MakeInt

// MakeInt provides a way of creating an ObservableInt from scratch by
// calling observer methods programmatically. A make function conforming to the
// MakeIntFunc signature will be called by MakeInt provining a Next, Error and
// Complete function that can be called by the code that implements the
// Observable.
func MakeInt(make MakeIntFunc) ObservableInt {
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Canceled() {
				return
			}
			next := func(n int) {
				if subscriber.Subscribed() {
					observe(n, nil, false)
				}
			}
			err := func(e error) {
				done = true
				if subscriber.Subscribed() {
					observe(zeroInt, e, true)
				}
			}
			complete := func() {
				done = true
				if subscriber.Subscribed() {
					observe(zeroInt, nil, true)
				}
			}
			make(next, err, complete)
			if !done && subscriber.Subscribed() {
				self()
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Schedulers

func TrampolineScheduler() Scheduler {
	return scheduler.Trampoline
}

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

//jig:name ObservableIntPrintln

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println is performed on the Trampoline scheduler.
func (o ObservableInt) Println() (err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next int, e error, done bool) {
		if !done {
			fmt.Println(next)
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
