// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ObserveOn

import (
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

//jig:name ObservableInt_ObserveOn

// ObserveOn specifies a dispatch function to use for delivering values to the observer.
func (o ObservableInt) ObserveOn(dispatch func(task func())) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			task := func() {
				observe(next, err, done)
			}
			dispatch(task)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:name ObservableInt_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Subscribe(observe IntObserver, schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	observer := func(next int, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero int
			observe(zero, err, true)
			subscriber.Done(err)
		}
	}
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	o(observer, schedulers[0], subscriber)
	return subscriber
}
