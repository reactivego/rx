// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package As

import (
	"errors"

	"github.com/reactivego/rx/schedulers"
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

//jig:name ObserveFunc

// ObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type ObserveFunc func(interface{}, error, bool)

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
			type observer_subscriber struct {
				ObserveFunc
				Subscriber
			}
			f(&observer_subscriber{observer, subscriber})
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

//jig:name StringObserveFunc

// StringObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type StringObserveFunc func(string, error, bool)

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
			type observer_subscriber struct {
				StringObserveFunc
				Subscriber
			}
			f(&observer_subscriber{observer, subscriber})
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

//jig:name Float64ObserveFunc

// Float64ObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type Float64ObserveFunc func(float64, error, bool)

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

//jig:name ErrTypecastToFloat64

// ErrTypecastToFloat64 is delivered to an observer if the generic value cannot be
// typecast to float64.
var ErrTypecastToFloat64 = errors.New("typecast to float64 failed")

//jig:name ObservableAsFloat64

// AsFloat64 turns an Observable of interface{} into an ObservableFloat64. If during
// observing a typecast fails, the error ErrTypecastToFloat64 will be emitted.
func (o Observable) AsFloat64() ObservableFloat64 {
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

//jig:name ObservableStringAsAny

// AsAny turns a typed ObservableString into an Observable of interface{}.
func (o ObservableString) AsAny() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next string, err error, done bool) {
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableFloat64

// ObservableFloat64 is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableFloat64 func(Float64ObserveFunc, Scheduler, Subscriber)

//jig:name NewScheduler

func NewGoroutine() Scheduler	{ return &schedulers.Goroutine{} }

func NewTrampoline() Scheduler	{ return &schedulers.Trampoline{} }

//jig:name SubscribeOptions

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

// SubscribeOptions is a struct with options for Subscribe related methods.
type SubscribeOptions struct {
	// SubscribeOn is the scheduler to run the observable subscription on.
	SubscribeOn	Scheduler
	// OnSubscribe is called right after the subscription is created and before
	// subscribing continues further.
	OnSubscribe	func(subscription Subscription)
	// OnUnsubscribe is called by the subscription to notify the client that the
	// subscription has been canceled.
	OnUnsubscribe	func()
}

// NewSubscriber will return a newly created subscriber. Before returning the
// subscription the OnSubscribe callback (if set) will already have been called.
func (options SubscribeOptions) NewSubscriber() Subscriber {
	subscription := subscriber.NewWithCallback(options.OnUnsubscribe)
	if options.OnSubscribe != nil {
		options.OnSubscribe(subscription)
	}
	return subscription
}

// SubscribeOptionSetter is the type of a function for setting SubscribeOptions.
type SubscribeOptionSetter func(options *SubscribeOptions)

// SubscribeOn takes the scheduler to run the observable subscription on and
// additional setters. It will first set the SubscribeOn option before
// calling the other setters provided as a parameter.
func SubscribeOn(subscribeOn Scheduler, setters ...SubscribeOptionSetter) SubscribeOptionSetter {
	return func(options *SubscribeOptions) {
		options.SubscribeOn = subscribeOn
		for _, setter := range setters {
			setter(options)
		}
	}
}

// OnSubscribe takes a callback to be called on subscription.
func OnSubscribe(callback func(Subscription)) SubscribeOptionSetter {
	return func(options *SubscribeOptions) { options.OnSubscribe = callback }
}

// OnUnsubscribe takes a callback to be called on subscription cancelation.
func OnUnsubscribe(callback func()) SubscribeOptionSetter {
	return func(options *SubscribeOptions) { options.OnUnsubscribe = callback }
}

// NewSubscribeOptions will create a new SubscribeOptions struct and then call
// the setter on it to recursively set all the options. It then returns a
// pointer to the created SubscribeOptions struct.
func NewSubscribeOptions(setter SubscribeOptionSetter) *SubscribeOptions {
	options := &SubscribeOptions{}
	setter(options)
	return options
}

//jig:name ObservableSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o Observable) Subscribe(observe ObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zero, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:name ObservableSubscribeNext

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscriber.
func (o Observable) SubscribeNext(f func(next interface{}), setters ...SubscribeOptionSetter) Subscription {
	return o.Subscribe(func(next interface{}, err error, done bool) {
		if !done {
			f(next)
		}
	}, setters...)
}

//jig:name ObservableFloat64Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableFloat64) Subscribe(observe Float64ObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next float64, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFloat64, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:name ObservableFloat64SubscribeNext

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscriber.
func (o ObservableFloat64) SubscribeNext(f func(next float64), setters ...SubscribeOptionSetter) Subscription {
	return o.Subscribe(func(next float64, err error, done bool) {
		if !done {
			f(next)
		}
	}, setters...)
}

//jig:name ObservableFloat64Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
func (o ObservableFloat64) Wait(setters ...SubscribeOptionSetter) (e error) {
	o.Subscribe(func(next float64, err error, done bool) {
		if done {
			e = err
		}
	}, setters...).Wait()
	return e
}
