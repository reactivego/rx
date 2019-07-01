// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package Reduce

import (
	"fmt"

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

//jig:name FromInt

// FromInt creates an ObservableInt from multiple int values passed in.
func FromInt(slice ...int) ObservableInt {
	return FromSliceInt(slice)
}

//jig:name ObservableIntReduce

// Reduce applies a reducer function to each item emitted by an ObservableInt
// and the previous reducer result. The operator accepts a seed argument that
// is passed to the reducer for the first item emitted by the ObservableInt.
// Reduce emits only the final value.
func (o ObservableInt) Reduce(reducer func(interface{}, int) interface{}, seed interface{}) Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next int, err error, done bool) {
			if !done {
				state = reducer(state, next)
			} else {
				if err == nil {
					observe(state, nil, false)
				}
				observe(zero, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntReduceFloat32

// ReduceFloat32 applies a reducer function to each item emitted by an ObservableInt
// and the previous reducer result. The operator accepts a seed argument that
// is passed to the reducer for the first item emitted by the ObservableInt.
// ReduceFloat32 emits only the final value.
func (o ObservableInt) ReduceFloat32(reducer func(float32, int) float32, seed float32) ObservableFloat32 {
	observable := func(observe Float32ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next int, err error, done bool) {
			if !done {
				state = reducer(state, next)
			} else {
				if err == nil {
					observe(state, nil, false)
				}
				observe(zeroFloat32, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
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

//jig:name Float32ObserveFunc

// Float32ObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type Float32ObserveFunc func(next float32, err error, done bool)

//jig:name zeroFloat32

var zeroFloat32 float32

//jig:name ObservableFloat32

// ObservableFloat32 is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableFloat32 func(Float32ObserveFunc, Scheduler, Subscriber)

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

//jig:name ObservableFloat32Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
func (o ObservableFloat32) Subscribe(observe Float32ObserveFunc, options ...SubscribeOption) Subscription {
	scheduler, subscriber := newSchedulerAndSubscriber(options)
	observer := func(next float32, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFloat32, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	return subscriber
}

//jig:name ObservableFloat32ToSingle

// ToSingle blocks until the ObservableFloat32 emits exactly one value or an error.
// The value and any error are returned.
//
// This function subscribes to the source observable on the NewGoroutine
// scheduler. The NewGoroutine scheduler works in more situations for
// complex chains of observables, like when merging the output of multiple
// observables.
func (o ObservableFloat32) ToSingle(options ...SubscribeOption) (entry float32, err error) {
	scheduler := NewGoroutineScheduler()
	o.Single().Subscribe(func(next float32, e error, done bool) {
		if !done {
			entry = next
		} else {
			err = e
		}
	}, SubscribeOn(scheduler, options...)).Wait()
	return
}

//jig:name ObservablePrintln

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
func (o Observable) Println() (err error) {
	subscriber := subscriber.New()
	scheduler := CurrentGoroutineScheduler()
	observer := func(next interface{}, e error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
}

//jig:name ObservableFloat32Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
func (o ObservableFloat32) Println() (err error) {
	subscriber := subscriber.New()
	scheduler := CurrentGoroutineScheduler()
	observer := func(next float32, e error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
}

//jig:name ObservableFloat32Single

// Single enforces that the observableFloat32 sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o ObservableFloat32) Single() ObservableFloat32 {
	return o.AsObservable().Single().AsObservableFloat32()
}

//jig:name ObservableFloat32AsObservable

// AsObservable turns a typed ObservableFloat32 into an Observable of interface{}.
func (o ObservableFloat32) AsObservable() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next float32, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name ObservableSingle

// Single enforces that the observable sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o Observable) Single() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			count	int
			latest	interface{}
		)
		observer := func(next interface{}, err error, done bool) {
			if count < 2 {
				if done {
					if err != nil {
						observe(nil, err, true)
					} else {
						if count == 1 {
							observe(latest, nil, false)
							observe(nil, nil, true)
						} else {
							observe(nil, RxError("expected one value, got none"), true)
						}
					}
				} else {
					count++
					if count == 1 {
						latest = next
					} else {
						observe(nil, RxError("expected one value, got multiple"), true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ErrTypecastToFloat32

// ErrTypecastToFloat32 is delivered to an observer if the generic value cannot be
// typecast to float32.
const ErrTypecastToFloat32 = RxError("typecast to float32 failed")

//jig:name ObservableAsObservableFloat32

// AsFloat32 turns an Observable of interface{} into an ObservableFloat32. If during
// observing a typecast fails, the error ErrTypecastToFloat32 will be emitted.
func (o Observable) AsObservableFloat32() ObservableFloat32 {
	observable := func(observe Float32ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextFloat32, ok := next.(float32); ok {
					observe(nextFloat32, err, done)
				} else {
					observe(zeroFloat32, ErrTypecastToFloat32, true)
				}
			} else {
				observe(zeroFloat32, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
