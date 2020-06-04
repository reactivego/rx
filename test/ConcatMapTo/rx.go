// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ConcatMapTo

import (
	"fmt"
	"sync"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
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
	var zeroInt int
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
					observe(zeroInt, nil, true)
				}
			}
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

//jig:name OfString

// OfString emits a variable amount of values in a sequence and then emits a
// complete notification.
func OfString(slice ...string) ObservableString {
	var zeroString string
	observable := func(observe StringObserver, scheduler Scheduler, subscriber Subscriber) {
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
					observe(zeroString, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableInt_ConcatMapToString

// ConcatMapToString maps every entry emitted by the ObservableInt into a single
// ObservableString. The stream of ObservableString items is then flattened by
// concattenating the emissions from the observables without interleaving.
func (o ObservableInt) ConcatMapToString(inner ObservableString) ObservableString {
	project := func(int) ObservableString { return inner }
	return o.MapObservableString(project).ConcatAll()
}

//jig:name ObservableInt_MapObservableString

// MapObservableString transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapObservableString(project func(int) ObservableString) ObservableObservableString {
	observable := func(observe ObservableStringObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			var mapped ObservableString
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
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

//jig:name ObservableObservableString_ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableString) ConcatAll() ObservableString {
	observable := func(observe StringObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var concat struct {
			sync.Mutex
			observables	[]ObservableString
			observer	StringObserver
			subscriber	Subscriber
		}

		var source struct {
			observer	ObservableStringObserver
			subscriber	Subscriber
		}

		concat.observer = func(next string, err error, done bool) {
			concat.Lock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(concat.observables) == 0 {
					if !source.subscriber.Subscribed() {
						var zero string
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

		source.observer = func(next ObservableString, err error, done bool) {
			if !done {
				concat.Lock()
				initial := concat.observables == nil
				concat.observables = append(concat.observables, next)
				concat.Unlock()
				if initial {
					var zero string
					concat.observer(zero, nil, true)
				}
			} else {
				concat.Lock()
				initial := concat.observables == nil
				source.subscriber.Done(err)
				concat.Unlock()
				if initial || err != nil {
					var zero string
					concat.observer(zero, err, true)
				}
			}
		}
		source.subscriber = subscriber.Add()
		o(source.observer, subscribeOn, source.subscriber)
	}
	return observable
}
