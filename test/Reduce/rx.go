// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Reduce

import (
	"fmt"

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

//jig:name ObservableInt_ReduceFloat32

// ReduceFloat32 applies a reducer function to each item emitted by an ObservableInt
// and the previous reducer result. The operator accepts a seed argument that
// is passed to the reducer for the first item emitted by the ObservableInt.
// ReduceFloat32 emits only the final value.
func (o ObservableInt) ReduceFloat32(reducer func(float32, int) float32, seed float32) ObservableFloat32 {
	observable := func(observe Float32Observer, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next int, err error, done bool) {
			if !done {
				state = reducer(state, next)
			} else {
				if err == nil {
					observe(state, nil, false)
				}
				var zeroFloat32 float32
				observe(zeroFloat32, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_Reduce

// Reduce applies a reducer function to each item emitted by an ObservableInt
// and the previous reducer result. The operator accepts a seed argument that
// is passed to the reducer for the first item emitted by the ObservableInt.
// Reduce emits only the final value.
func (o ObservableInt) Reduce(reducer func(interface{}, int) interface{}, seed interface{}) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next int, err error, done bool) {
			if !done {
				state = reducer(state, next)
			} else {
				if err == nil {
					observe(state, nil, false)
				}
				var zero interface{}
				observe(zero, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Float32Observer

// Float32Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Float32Observer func(next float32, err error, done bool)

//jig:name ObservableFloat32

// ObservableFloat32 is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableFloat32 func(Float32Observer, Scheduler, Subscriber)

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

//jig:name ObservableFloat32_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFloat32) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next float32, err error, done bool) {
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
