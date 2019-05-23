// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package SwitchMap

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/rx/schedulers"
	"github.com/reactivego/rx/subscriber"
)

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
			type ObserverSubscriber struct {
				ObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name Never

// Never creates an Observable that emits no items and does't terminate.
func Never() Observable {
	return Create(func(observer Observer) {
	})
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

//jig:name IntObserveFunc

// IntObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type IntObserveFunc func(int, error, bool)

var zeroInt int

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

//jig:name ObservableInt

// ObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableInt func(IntObserveFunc, Scheduler, Subscriber)

//jig:name IntObserver

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

//jig:name Interval

// Interval creates an ObservableInt that emits a sequence of integers spaced
// by a particular time interval.
func Interval(interval time.Duration) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for i := 0; ; i++ {
			time.Sleep(interval)
			if observer.Closed() {
				return
			}
			observer.Next(i)
		}
	})
}

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

//jig:name EmptyString

// EmptyString creates an Observable that emits no items but terminates normally.
func EmptyString() ObservableString {
	return CreateString(func(observer StringObserver) {
		observer.Complete()
	})
}

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler interface {
	Schedule(task func())
}

//jig:name Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

//jig:name ConstError

type Error string

func (e Error) Error() string	{ return string(e) }

//jig:name ObservableSerialize

// Serialize forces an Observable to make serialized calls and to be
// well-behaved.
func (o Observable) Serialize() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex		sync.Mutex
			alreadyDone	bool
		)
		observer := func(next interface{}, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !alreadyDone {
				alreadyDone = done
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableTimeout

// ErrTimeout is delivered to an observer if the stream times out.
const ErrTimeout = Error("timeout")

// Timeout mirrors the source Observable, but issues an error notification if a
// particular period of time elapses without any emitted items.
//
// This observer starts a goroutine for every subscription to monitor the
// timeout deadline. It is guaranteed that calls to the observer for this
// subscription will never be called concurrently. It is however almost certain
// that any timeout error will be delivered on a goroutine other than the one
// delivering the next values.
func (o Observable) Timeout(timeout time.Duration) Observable {
	observable := Observable(func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		deadline := time.NewTimer(timeout)
		unsubscribe := make(chan struct{})
		observer := func(next interface{}, err error, done bool) {
			if deadline.Stop() {
				if subscriber.Closed() {
					return
				}
				observe(next, err, done)
				if done {
					return
				}
				deadline.Reset(timeout)
			}
		}
		watchdog := func() {
			select {
			case <-deadline.C:
				if subscriber.Closed() {
					return
				}
				observe(nil, ErrTimeout, true)
			case <-unsubscribe:
			}
		}
		go watchdog()
		o(observer, subscribeOn, subscriber.Add(func() { close(unsubscribe) }))
	})
	return observable.Serialize()
}

//jig:name ErrTypecastToString

// ErrTypecastToString is delivered to an observer if the generic value cannot be
// typecast to string.
const ErrTypecastToString = Error("typecast to string failed")

//jig:name ObservableAsObservableString

// AsString turns an Observable of interface{} into an ObservableString. If during
// observing a typecast fails, the error ErrTypecastToString will be emitted.
func (o Observable) AsObservableString() ObservableString {
	observable := func(observe StringObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextString, ok := next.(string); ok {
					observe(nextString, err, done)
				} else {
					observe(zeroString, ErrTypecastToString, true)
				}
			} else {
				observe(zeroString, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableTake

// Take emits only the first n items emitted by an Observable.
func (o Observable) Take(n int) Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		taken := 0
		observer := func(next interface{}, err error, done bool) {
			if taken < n {
				observe(next, err, done)
				if !done {
					taken++
					if taken >= n {
						observe(nil, nil, true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntTake

// Take emits only the first n items emitted by an ObservableInt.
func (o ObservableInt) Take(n int) ObservableInt {
	return o.AsObservable().Take(n).AsObservableInt()
}

//jig:name ObservableCatch

// Catch recovers from an error notification by continuing the sequence without
// emitting the error but by switching to the catch Observable to provide
// items.
func (o Observable) Catch(catch Observable) Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if err != nil {
				catch(observe, subscribeOn, subscriber)
			} else {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntSwitchMapString

// SwitchMapString transforms the items emitted by an ObservableInt by applying a
// function to each item an returning an ObservableString. In doing so, it behaves much like
// MergeMap (previously FlatMap), except that whenever a new ObservableString is emitted
// SwitchMap will unsubscribe from the previous ObservableString and begin emitting items
// from the newly emitted one.
func (o ObservableInt) SwitchMapString(project func(int) ObservableString) ObservableString {
	return o.MapObservableString(project).SwitchAll()
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

//jig:name ObservableStringSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableString) Subscribe(observe StringObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next string, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroString, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:name ObservableStringToSlice

// ToSlice collects all values from the ObservableString into an slice. The
// complete slice and any error are returned.
//
// This function subscribes to the source observable on the Goroutine scheduler.
// The Goroutine scheduler works in more situations for complex chains of
// observables, like when merging the output of multiple observables.
func (o ObservableString) ToSlice(setters ...SubscribeOptionSetter) (a []string, e error) {
	scheduler := NewGoroutine()
	o.Subscribe(func(next string, err error, done bool) {
		if !done {
			a = append(a, next)
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, setters...)).Wait()
	return a, e
}

//jig:name ObservableIntMapObservableString

// MapObservableString transforms the items emitted by an ObservableInt by applying a
// function to each item.
func (o ObservableInt) MapObservableString(project func(int) ObservableString) ObservableObservableString {
	observable := func(observe ObservableStringObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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

//jig:name ErrTypecastToInt

// ErrTypecastToInt is delivered to an observer if the generic value cannot be
// typecast to int.
const ErrTypecastToInt = Error("typecast to int failed")

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

//jig:name ObservableStringObserveFunc

// ObservableStringObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type ObservableStringObserveFunc func(ObservableString, error, bool)

var zeroObservableString ObservableString

// Next is called by an ObservableObservableString to emit the next ObservableString value to the
// observer.
func (f ObservableStringObserveFunc) Next(next ObservableString) {
	f(next, nil, false)
}

// Error is called by an ObservableObservableString to report an error to the observer.
func (f ObservableStringObserveFunc) Error(err error) {
	f(zeroObservableString, err, true)
}

// Complete is called by an ObservableObservableString to signal that no more data is
// forthcoming to the observer.
func (f ObservableStringObserveFunc) Complete() {
	f(zeroObservableString, nil, true)
}

//jig:name ObservableObservableString

// ObservableObservableString is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableObservableString func(ObservableStringObserveFunc, Scheduler, Subscriber)

//jig:name LinkEnums

// state
const (
	linkUnsubscribed	= iota
	linkSubscribing
	linkIdle
	linkBusy
	linkError	// done:error
	linkCanceled	// externally:canceled
	linkCompleting
	linkComplete	// done:complete
)

// callbackState
const (
	callbackNil	= iota
	settingCallback
	callbackSet
)

// callbackKind
const (
	linkCallbackOnComplete	= iota
	linkCancelOrCompleted
)

//jig:name StringLink

type StringLinkObserveFunc func(*StringLink, string, error, bool)

type StringLink struct {
	observe		StringLinkObserveFunc
	state		int32
	callbackState	int32
	callbackKind	int
	callback	func()
	subscriber	Subscriber
}

func NewInitialStringLink() *StringLink {
	return &StringLink{state: linkCompleting, subscriber: subscriber.New()}
}

func NewStringLink(observe StringLinkObserveFunc, subscriber Subscriber) *StringLink {
	return &StringLink{
		observe:	observe,
		subscriber:	subscriber.Add(func() {}),
	}
}

func (o *StringLink) Observe(next string, err error, done bool) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkIdle, linkBusy) {
		if atomic.LoadInt32(&o.state) > linkBusy {
			return Error("Already Done")
		}
		return Error("Recursion Error")
	}
	o.observe(o, next, err, done)
	if done {
		if err != nil {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkError) {
				return Error("Internal Error: 'busy' -> 'error'")
			}
		} else {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkCompleting) {
				return Error("Internal Error: 'busy' -> 'completing'")
			}
		}
	} else {
		if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkIdle) {
			return Error("Internal Error: 'busy' -> 'idle'")
		}
	}
	if atomic.LoadInt32(&o.callbackState) != callbackSet {
		return nil
	}
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	if o.callbackKind == linkCancelOrCompleted {
		if atomic.CompareAndSwapInt32(&o.state, linkIdle, linkCanceled) {
			o.callback()
		}
	}
	return nil
}

func (o *StringLink) SubscribeTo(observable ObservableString, scheduler Scheduler) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkUnsubscribed, linkSubscribing) {
		return Error("Already Subscribed")
	}
	observer := func(next string, err error, done bool) {
		o.Observe(next, err, done)
	}
	observable(observer, scheduler, o.subscriber)
	if !atomic.CompareAndSwapInt32(&o.state, linkSubscribing, linkIdle) {
		return Error("Internal Error")
	}
	return nil
}

func (o *StringLink) Cancel(callback func()) error {
	if !atomic.CompareAndSwapInt32(&o.callbackState, callbackNil, settingCallback) {
		return Error("Already Waiting")
	}
	o.callbackKind = linkCancelOrCompleted
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return Error("Internal Error")
	}
	o.subscriber.Unsubscribe()
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	if atomic.CompareAndSwapInt32(&o.state, linkIdle, linkCanceled) {
		o.callback()
	}
	return nil
}

func (o *StringLink) OnComplete(callback func()) error {
	if !atomic.CompareAndSwapInt32(&o.callbackState, callbackNil, settingCallback) {
		return Error("Already Waiting")
	}
	o.callbackKind = linkCallbackOnComplete
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return Error("Internal Error")
	}
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	return nil
}

//jig:name ObservableObservableStringSwitchAll

// SwitchAll converts an Observable that emits Observables into a single Observable
// that emits the items emitted by the most-recently-emitted of those Observables.
func (o ObservableObservableString) SwitchAll() ObservableString {
	observable := func(observe StringObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(link *StringLink, next string, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				link.subscriber.Unsubscribe()
			}
		}
		currentLink := NewInitialStringLink()
		var switcherMutex sync.Mutex
		switcherSubscriber := subscriber.Add(func() {})
		switcher := func(next ObservableString, err error, done bool) {
			switch {
			case !done:
				previousLink := currentLink
				func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink = NewStringLink(observer, subscriber)
				}()
				previousLink.Cancel(func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink.SubscribeTo(next, subscribeOn)
				})
			case err != nil:
				currentLink.Cancel(func() {
					observe(zeroString, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					observe(zeroString, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}