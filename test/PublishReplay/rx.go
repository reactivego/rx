// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package PublishReplay

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/reactivego/channel"
	"github.com/reactivego/rx/schedulers"
	"github.com/reactivego/rx/subscriber"
)

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

//jig:name FromChanInt

// FromChanInt creates an ObservableInt from a Go channel of int values.
// It's not possible for the code feeding into the channel to send an error.
// The feeding code can send zero or more int items and then closing the
// channel will be seen as completion.
func FromChanInt(ch <-chan int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for next := range ch {
			if observer.Closed() {
				return
			}
			observer.Next(next)
		}
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

//jig:name ObservableIntSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableInt) Subscribe(observe IntObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next int, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroInt, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:name ConnectableInt

// ConnectableInt is an ObservableInt that has an additional method Connect()
// used to Subscribe to the parent observable and then multicasting values to
// all subscribers of ConnectableInt.
type ConnectableInt struct {
	ObservableInt
	connect	func(options []SubscribeOptionSetter) Subscription
}

//jig:name ObservableIntMulticast

// Multicast converts an ordinary Observable into a connectable Observable.
// A connectable observable will only start emitting values after its Connect
// method has been called. The factory method passed in should return a
// new SubjectInt that implements the actual multicasting behavior.
func (o ObservableInt) Multicast(factory func() SubjectInt) ConnectableInt {
	const (
		active	int32	= iota
		notifying
		terminated
	)
	var subject struct {
		state	int32
		atomic.Value
	}
	subject.Store(factory())
	observable := func(observe IntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if s, ok := subject.Load().(SubjectInt); ok {
			s.ObservableInt(observe, subscribeOn, subscriber)
		}
	}
	observer := func(next int, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subject.state, active, notifying) {
			if s, ok := subject.Load().(SubjectInt); ok {
				s.IntObserveFunc(next, err, done)
			}
			if !done {
				atomic.CompareAndSwapInt32(&subject.state, notifying, active)
			} else {
				atomic.CompareAndSwapInt32(&subject.state, notifying, terminated)
			}
		}
	}
	const (
		unsubscribed	int32	= iota
		subscribed
	)
	var subscriber struct {
		state	int32
		atomic.Value
	}
	connect := func(setters []SubscribeOptionSetter) Subscription {
		if atomic.CompareAndSwapInt32(&subject.state, terminated, active) {
			subject.Store(factory())
		}
		if atomic.CompareAndSwapInt32(&subscriber.state, unsubscribed, subscribed) {
			scheduler := NewGoroutine()
			setter := SubscribeOn(scheduler, setters...)
			subscription := o.Subscribe(observer, setter)
			subscriber.Store(subscription)
			subscription.Add(func() {
				atomic.CompareAndSwapInt32(&subscriber.state, subscribed, unsubscribed)
			})
		}
		subscription := subscriber.Load().(Subscriber)
		return subscription.Add(func() { subscription.Unsubscribe() })
	}
	return ConnectableInt{ObservableInt: observable, connect: connect}
}

//jig:name SubjectInt

// SubjectInt is a combination of an observer and observable. Subjects are
// special because they are the only reactive constructs that support
// multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// A SubjectInt embeds ObservableInt and IntObserveFunc. This exposes the
// methods and fields of both types on SubjectInt. Use the ObservableInt
// methods to subscribe to it. Use the IntObserveFunc Next, Error and Complete
// methods to feed data to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubjectInt
// functions for more info.
//
// Important! a subject is a hot observable. This means that subscribing to
// it will block the calling goroutine while it is waiting for items and
// notifications to receive. Unless you have code on a different goroutine
// already feeding into the subject, your subscribe will deadlock.
// Alternatively, you could subscribe on a goroutine as shown in the example.
type SubjectInt struct {
	ObservableInt
	IntObserveFunc
}

//jig:name MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:name NewReplaySubjectInt

// NewReplaySubjectInt creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
//
// Note that this implementation is non-blocking. When no subscribers are
// present the buffer fills up to bufferCapacity after which new items will
// start overwriting the oldest ones according to the FIFO principle.
// If a subscriber cannot keep up with the data rate of the source observable,
// eventually the buffer for the subscriber will overflow. At that moment the
// subscriber will receive an ErrMissingBackpressure error.
func NewReplaySubjectInt(bufferCapacity int, windowDuration time.Duration) SubjectInt {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	ch := channel.NewChan(bufferCapacity, 16)

	observable := Observable(func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		ep, err := ch.NewEndpoint(channel.ReplayAll)
		if err != nil {
			observe(nil, err, true)
			return
		}
		observable := Create(func(observer Observer) {
			receive := func(value interface{}, err error, closed bool) bool {
				if !closed {
					observer.Next(value)
				} else {
					observer.Error(err)
				}
				return !observer.Closed()
			}
			ep.Range(receive, windowDuration)
		})
		observable(observe, subscribeOn, subscriber.Add(ep.Cancel))
	})

	observer := func(next int, err error, done bool) {
		if !ch.Closed() {
			if !done {
				ch.Send(next)
			} else {
				ch.Close(err)
			}
		}
	}

	return SubjectInt{observable.AsObservableInt(), observer}
}

//jig:name ObservableIntPublishReplay

// Replay uses Multicast to control the subscription of a ReplaySubject to a
// source observable and turns the subject into a connectable observable.
// A ReplaySubject emits to any observer all of the items that were emitted by
// the source observable, regardless of when the observer subscribes.
//
// If the source completed and as a result the internal ReplaySubject
// terminated, then calling Connect again will replace the old ReplaySubject
// with a newly created one.
func (o ObservableInt) PublishReplay(bufferCapacity int, windowDuration time.Duration) ConnectableInt {
	factory := func() SubjectInt {
		return NewReplaySubjectInt(bufferCapacity, windowDuration)
	}
	return o.Multicast(factory)
}

//jig:name ConnectableIntConnect

// Connect instructs a connectable Observable to begin emitting items to its
// subscribers. All values will then be passed on to the observers that
// subscribed to this connectable observable
func (c ConnectableInt) Connect(setters ...SubscribeOptionSetter) Subscription {
	return c.connect(setters)
}

//jig:name ObservableIntWait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
func (o ObservableInt) Wait(setters ...SubscribeOptionSetter) (e error) {
	o.Subscribe(func(next int, err error, done bool) {
		if done {
			e = err
		}
	}, setters...).Wait()
	return e
}

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

//jig:name ObservableIntToSlice

// ToSlice collects all values from the ObservableInt into an slice. The
// complete slice and any error are returned.
//
// This function subscribes to the source observable on the Goroutine scheduler.
// The Goroutine scheduler works in more situations for complex chains of
// observables, like when merging the output of multiple observables.
func (o ObservableInt) ToSlice(setters ...SubscribeOptionSetter) (a []int, e error) {
	scheduler := NewGoroutine()
	o.Subscribe(func(next int, err error, done bool) {
		if !done {
			a = append(a, next)
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, setters...)).Wait()
	return a, e
}

//jig:name ConnectableIntAutoConnect

// AutoConnect makes a ConnectableInt behave like an ordinary ObservableInt that
// automatically connects when the specified number of clients subscribe to it.
// If count is 0, then AutoConnect will immediately call connect on the
// ConnectableInt before returning the ObservableInt part of the ConnectableInt.
func (o ConnectableInt) AutoConnect(count int, setters ...SubscribeOptionSetter) ObservableInt {
	if count == 0 {
		o.connect(setters)
		return o.ObservableInt
	}
	var refcount int32
	observable := func(observe IntObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&refcount, 1) == int32(count) {
			o.connect(setters)
		}
		o.ObservableInt(observe, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ErrTypecastToInt

// ErrTypecastToInt is delivered to an observer if the generic value cannot be
// typecast to int.
var ErrTypecastToInt = errors.New("typecast to int failed")

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
