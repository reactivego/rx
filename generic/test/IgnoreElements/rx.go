// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package IgnoreElements

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
	observable := func(observe IntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		subscribeOn.ScheduleRecursive(func(self func()) {
			if !subscriber.Canceled() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if !subscriber.Canceled() {
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

//jig:name FromInts

// FromInts creates an ObservableInt from multiple int values passed in.
func FromInts(slice ...int) ObservableInt {
	return FromSliceInt(slice)
}

//jig:name ObservableIgnoreElements

// IgnoreElements does not emit any items from an Observable but mirrors its termination notification.
func (o Observable) IgnoreElements() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if done {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntIgnoreElements

// IgnoreElements does not emit any items from an ObservableInt but mirrors its termination notification.
func (o ObservableInt) IgnoreElements() ObservableInt {
	return o.AsObservable().IgnoreElements().AsObservableInt()
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
// This function subscribes to the source observable on the NewGoroutine
// scheduler. The NewGoroutine scheduler works in more situations for
// complex chains of observables, like when merging the output of multiple
// observables.
func (o ObservableInt) ToSlice(options ...SubscribeOption) (slice []int, err error) {
	scheduler := NewGoroutineScheduler()
	o.Subscribe(func(next int, e error, done bool) {
		if !done {
			slice = append(slice, next)
		} else {
			err = e
		}
	}, SubscribeOn(scheduler, options...)).Wait()
	return
}

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
