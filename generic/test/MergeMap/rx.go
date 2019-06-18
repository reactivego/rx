// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package MergeMap

import (
	"sync"
	"sync/atomic"

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

//jig:name IntObserveFunc

// IntObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type IntObserveFunc func(next int, err error, done bool)

var zeroInt int

//jig:name ObservableInt

// ObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableInt func(IntObserveFunc, Scheduler, Subscriber)

//jig:name IntObserver

// Next is called by an ObservableInt to emit the next int value to the
// observer.
func (f IntObserveFunc) Next(next int) {
	f(next, nil, false)
}

// Error is called by an ObservableInt to report an error to the observer.
func (f IntObserveFunc) Error(err error) {
	f(zeroInt, err, true)
}

// Complete is called by an ObservableInt to signal that no more data is
// forthcoming to the observer.
func (f IntObserveFunc) Complete() {
	f(zeroInt, nil, true)
}

// IntObserver is the interface used with CreateInt when implementing a custom
// observable.
type IntObserver interface {
	// Next emits the next int value.
	Next(int)
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:name CreateInt

// CreateInt creates an Observable from scratch by calling observer methods
// programmatically.
func CreateInt(f func(IntObserver)) ObservableInt {
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next int, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				IntObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name Range

// Range creates an ObservableInt that emits a range of sequential integers.
func Range(start, count int) ObservableInt {
	end := start + count
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := start
		scheduler.ScheduleRecursive(func(self func()) {
			if !subscriber.Closed() {
				if i < end {
					observe(i, nil, false)
					if !subscriber.Closed() {
						i++
						self()
					}
				} else {
					observe(zeroInt, nil, true)
				}
			}
		})
	}
	return observable
}

//jig:name ObservableIntMergeMapInt

// MergeMapInt transforms the items emitted by an ObservableInt by applying a
// function to each item an returning an ObservableInt. The stream of ObservableInt
// items is then merged into a single stream of Int items using the MergeAll operator.
func (o ObservableInt) MergeMapInt(project func(int) ObservableInt) ObservableInt {
	return o.MapObservableInt(project).MergeAll()
}

//jig:name Schedulers

func ImmediateScheduler() Scheduler	{ return scheduler.Immediate }

func CurrentGoroutineScheduler() Scheduler	{ return scheduler.CurrentGoroutine }

func NewGoroutineScheduler() Scheduler	{ return scheduler.NewGoroutine }

//jig:name SubscribeOption

// SubscribeOption is an option that can be passed to the Subscribe method.
type SubscribeOption func(options *subscribeOptions)

type subscribeOptions struct {
	scheduler	Scheduler
	subscriber	Subscriber
	onSubscribe	func(subscription Subscription)
	onUnsubscribe	func()
}

// SubscribeOn returns an option that can be passed to the Subscribe method.
// It takes the scheduler to subscribe the observable on. The tasks that
// actually perform the observable functionality are scheduled on this
// scheduler. The other options that can be passed here are applied after the
// scheduler was set so any schedulers passed in via other will override
// the scheduler passed here.
func SubscribeOn(scheduler Scheduler, other ...SubscribeOption) SubscribeOption {
	return func(options *subscribeOptions) {
		options.scheduler = scheduler
		for _, setter := range other {
			setter(options)
		}
	}
}

// WithSubscriber returns an option that can be passed to the Subscribe method.
// The Subscribe method will use the subscriber passed here instead of creating
// a new one.
func WithSubscriber(subscriber Subscriber) SubscribeOption {
	return func(options *subscribeOptions) {
		options.subscriber = subscriber
	}
}

// OnSubscribe returns an option that can be passed to the Subscribe method.
// It takes a callback that is called from the Subscribe method just before
// subscribing continues further.
func OnSubscribe(callback func(Subscription)) SubscribeOption {
	return func(options *subscribeOptions) { options.onSubscribe = callback }
}

// OnUnsubscribe returns an option that can be passed to the Subscribe method.
// It takes a callback that is called by the Subscribe method to notify the
// client that the subscription has been canceled.
func OnUnsubscribe(callback func()) SubscribeOption {
	return func(options *subscribeOptions) { options.onUnsubscribe = callback }
}

// newSchedulerAndSubscriber will return either return the scheduler and subscriber
// passed in through the SubscribeOn() and WithSubscriber() options or it will
// return newly created scheduler and subscriber. Before returning the callback
// passed in through OnSubscribe() will already have been called.
func newSchedulerAndSubscriber(setters []SubscribeOption) (Scheduler, Subscriber) {
	options := &subscribeOptions{scheduler: CurrentGoroutineScheduler()}
	for _, setter := range setters {
		setter(options)
	}
	if options.subscriber == nil {
		options.subscriber = subscriber.New()
	}
	options.subscriber.OnUnsubscribe(options.onUnsubscribe)
	if options.onSubscribe != nil {
		options.onSubscribe(options.subscriber)
	}
	return options.scheduler, options.subscriber
}

//jig:name ObservableIntSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
func (o ObservableInt) Subscribe(observe IntObserveFunc, options ...SubscribeOption) Subscription {
	scheduler, subscriber := newSchedulerAndSubscriber(options)
	observer := func(next int, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroInt, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	return subscriber
}

//jig:name ObservableIntToSlice

// ToSlice collects all values from the ObservableInt into an slice. The
// complete slice and any error are returned.
//
// This function subscribes to the source observable on the Goroutine scheduler.
// The Goroutine scheduler works in more situations for complex chains of
// observables, like when merging the output of multiple observables.
func (o ObservableInt) ToSlice(options ...SubscribeOption) (a []int, e error) {
	scheduler := NewGoroutineScheduler()
	o.Subscribe(func(next int, err error, done bool) {
		if !done {
			a = append(a, next)
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, options...)).Wait()
	return a, e
}

//jig:name ObservableIntMapObservableInt

// MapObservableInt transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapObservableInt(project func(int) ObservableInt) ObservableObservableInt {
	observable := func(observe ObservableIntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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

//jig:name ObservableIntObserveFunc

// ObservableIntObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type ObservableIntObserveFunc func(next ObservableInt, err error, done bool)

var zeroObservableInt ObservableInt

//jig:name ObservableObservableInt

// ObservableObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableObservableInt func(ObservableIntObserveFunc, Scheduler, Subscriber)

//jig:name ObservableObservableIntMergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableInt) MergeAll() ObservableInt {
	observable := func(observe IntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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
					observe(zeroInt, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						observe(zeroInt, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableInt, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				observer(zeroInt, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
				observers.len = 1
				o(merger, subscribeOn, subscriber)
			}
		})
	}
	return observable
}
