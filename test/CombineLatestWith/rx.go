// Code generated by jig; DO NOT EDIT.

//go:generate jig

package CombineLatestWith

import (
	"fmt"
	"sync"

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

//jig:name From

// From creates an Observable from multiple interface{} values passed in.
func From(slice ...interface{}) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
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
					var zero interface{}
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Slice

type Slice = []interface{}

//jig:name ObservableObservable_CombineLatestAll

// CombineLatestAll flattens a higher order observable
// (e.g. ObservableObservable) by subscribing to
// all emitted observables (ie. Observable entries) until the source
// completes. It will then wait for all of the subscribed Observables
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableObservable) CombineLatestAll() ObservableSlice {
	observable := func(observe SliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observables := []Observable(nil)
		var observers struct {
			sync.Mutex
			assigned	[]bool
			values		[]interface{}
			initialized	int
			active		int
		}
		makeObserver := func(index int) Observer {
			observer := func(next interface{}, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if observers.active > 0 {
					switch {
					case !done:
						if !observers.assigned[index] {
							observers.assigned[index] = true
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						var zero []interface{}
						observe(zero, err, true)
					default:
						if observers.active--; observers.active == 0 {
							var zero []interface{}
							observe(zero, nil, true)
						}
					}
				}
			}
			return observer
		}

		observer := func(next Observable, err error, done bool) {
			switch {
			case !done:
				observables = append(observables, next)
			case err != nil:
				var zero []interface{}
				observe(zero, err, true)
			default:
				subscribeOn.Schedule(func() {
					if subscriber.Subscribed() {
						numObservables := len(observables)
						observers.assigned = make([]bool, numObservables)
						observers.values = make([]interface{}, numObservables)
						observers.active = numObservables
						for i, v := range observables {
							if !subscriber.Subscribed() {
								return
							}
							v(makeObserver(i), subscribeOn, subscriber)
						}
					}
				})
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_CombineLatestWith

// CombineLatestWith will subscribe to its Observable and all other
// Observables passed in. It will then wait for all of the ObservableBars
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o Observable) CombineLatestWith(other ...Observable) ObservableSlice {
	return FromObservable(append([]Observable{o}, other...)...).CombineLatestAll()
}

//jig:name ObservableObserver

// ObservableObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableObserver func(next Observable, err error, done bool)

//jig:name ObservableObservable

// ObservableObservable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservable func(ObservableObserver, Scheduler, Subscriber)

//jig:name SliceObserver

// SliceObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type SliceObserver func(next Slice, err error, done bool)

//jig:name ObservableSlice

// ObservableSlice is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableSlice func(SliceObserver, Scheduler, Subscriber)

//jig:name FromObservable

// FromObservable creates an ObservableObservable from multiple Observable values passed in.
func FromObservable(slice ...Observable) ObservableObservable {
	observable := func(observe ObservableObserver, scheduler Scheduler, subscriber Subscriber) {
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
					var zero Observable
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableSlice_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableSlice) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next Slice, err error, done bool) {
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
