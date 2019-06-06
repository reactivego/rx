// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package AsObservable

import (
	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler interface {
	Schedule(task func())
}

//jig:name Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

//jig:name StringObserveFunc

// StringObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type StringObserveFunc func(next string, err error, done bool)

var zeroString string

// Next is called by an ObservableString to emit the next string value to the
// observer.
func (f StringObserveFunc) Next(next string) {
	f(next, nil, false)
}

// Error is called by an ObservableString to report an error to the observer.
func (f StringObserveFunc) Error(err error) {
	f(zeroString, err, true)
}

// Complete is called by an ObservableString to signal that no more data is
// forthcoming to the observer.
func (f StringObserveFunc) Complete() {
	f(zeroString, nil, true)
}

//jig:name ObservableString

// ObservableString is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableString func(StringObserveFunc, Scheduler, Subscriber)

//jig:name StringObserver

// StringObserver is the interface used with CreateString when implementing a custom
// observable.
type StringObserver interface {
	// Next emits the next string value.
	Next(string)
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:name CreateString

// CreateString creates an Observable from scratch by calling observer methods
// programmatically.
func CreateString(f func(StringObserver)) ObservableString {
	observable := func(observe StringObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next string, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				StringObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name FromSliceString

// FromSliceString creates an ObservableString from a slice of string values passed in.
func FromSliceString(slice []string) ObservableString {
	return CreateString(func(observer StringObserver) {
		for _, next := range slice {
			if observer.Closed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

//jig:name FromString

// FromString creates an ObservableString from multiple string values passed in.
func FromString(slice ...string) ObservableString {
	return FromSliceString(slice)
}

//jig:name ObserveFunc

// ObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type ObserveFunc func(next interface{}, err error, done bool)

var zero interface{}

// Next is called by an Observable to emit the next interface{} value to the
// observer.
func (f ObserveFunc) Next(next interface{}) {
	f(next, nil, false)
}

// Error is called by an Observable to report an error to the observer.
func (f ObserveFunc) Error(err error) {
	f(zero, err, true)
}

// Complete is called by an Observable to signal that no more data is
// forthcoming to the observer.
func (f ObserveFunc) Complete() {
	f(zero, nil, true)
}

//jig:name Observable

// Observable is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type Observable func(ObserveFunc, Scheduler, Subscriber)

//jig:name Observer

// Observer is the interface used with Create when implementing a custom
// observable.
type Observer interface {
	// Next emits the next interface{} value.
	Next(interface{})
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:name Create

// Create creates an Observable from scratch by calling observer methods
// programmatically.
func Create(f func(Observer)) Observable {
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next interface{}, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				ObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name FromSlice

// FromSlice creates an Observable from a slice of interface{} values passed in.
func FromSlice(slice []interface{}) Observable {
	return Create(func(observer Observer) {
		for _, next := range slice {
			if observer.Closed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

//jig:name From

// From creates an Observable from multiple interface{} values passed in.
func From(slice ...interface{}) Observable {
	return FromSlice(slice)
}

//jig:name ObservableStringAsObservable

// AsObservable turns a typed ObservableString into an Observable of interface{}.
func (o ObservableString) AsObservable() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next string, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Float64ObserveFunc

// Float64ObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type Float64ObserveFunc func(next float64, err error, done bool)

var zeroFloat64 float64

// Next is called by an ObservableFloat64 to emit the next float64 value to the
// observer.
func (f Float64ObserveFunc) Next(next float64) {
	f(next, nil, false)
}

// Error is called by an ObservableFloat64 to report an error to the observer.
func (f Float64ObserveFunc) Error(err error) {
	f(zeroFloat64, err, true)
}

// Complete is called by an ObservableFloat64 to signal that no more data is
// forthcoming to the observer.
func (f Float64ObserveFunc) Complete() {
	f(zeroFloat64, nil, true)
}

//jig:name ObservableFloat64

// ObservableFloat64 is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableFloat64 func(Float64ObserveFunc, Scheduler, Subscriber)

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name ErrTypecastToFloat64

// ErrTypecastToFloat64 is delivered to an observer if the generic value cannot be
// typecast to float64.
const ErrTypecastToFloat64 = RxError("typecast to float64 failed")

//jig:name ObservableAsObservableFloat64

// AsFloat64 turns an Observable of interface{} into an ObservableFloat64. If during
// observing a typecast fails, the error ErrTypecastToFloat64 will be emitted.
func (o Observable) AsObservableFloat64() ObservableFloat64 {
	observable := func(observe Float64ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextFloat64, ok := next.(float64); ok {
					observe(nextFloat64, err, done)
				} else {
					observe(zeroFloat64, ErrTypecastToFloat64, true)
				}
			} else {
				observe(zeroFloat64, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name NewScheduler

func NewGoroutineScheduler() Scheduler	{ return &scheduler.Goroutine{} }

func NewTrampolineScheduler() Scheduler	{ return &scheduler.Trampoline{} }

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
	options := &subscribeOptions{scheduler: NewTrampolineScheduler()}
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

//jig:name ObservableSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
func (o Observable) Subscribe(observe ObserveFunc, options ...SubscribeOption) Subscription {
	scheduler, subscriber := newSchedulerAndSubscriber(options)
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zero, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	return subscriber
}

//jig:name ObservableSubscribeNext

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscription.
func (o Observable) SubscribeNext(f func(next interface{}), options ...SubscribeOption) Subscription {
	return o.Subscribe(func(next interface{}, err error, done bool) {
		if !done {
			f(next)
		}
	}, options...)
}

//jig:name ObservableFloat64Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
func (o ObservableFloat64) Subscribe(observe Float64ObserveFunc, options ...SubscribeOption) Subscription {
	scheduler, subscriber := newSchedulerAndSubscriber(options)
	observer := func(next float64, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFloat64, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	return subscriber
}

//jig:name ObservableFloat64SubscribeNext

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscription.
func (o ObservableFloat64) SubscribeNext(f func(next float64), options ...SubscribeOption) Subscription {
	return o.Subscribe(func(next float64, err error, done bool) {
		if !done {
			f(next)
		}
	}, options...)
}

//jig:name ObservableFloat64Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
func (o ObservableFloat64) Wait(options ...SubscribeOption) (e error) {
	o.Subscribe(func(next float64, err error, done bool) {
		if done {
			e = err
		}
	}, options...).Wait()
	return e
}
