// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Subject

import (
	"github.com/reactivego/multicast"
	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

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

//jig:name IntObserveFuncMethods

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

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

// NewSubscriber creates a new subscriber.
func NewSubscriber() Subscriber {
	return subscriber.New()
}

//jig:name ObservableInt

// ObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableInt func(IntObserveFunc, Scheduler, Subscriber)

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
type SubjectInt struct {
	ObservableInt
	IntObserveFunc
}

//jig:name NewSubjectInt

// NewSubjectInt creates a new Subject. After the subject is
// terminated, all subsequent subscriptions to the observable side will be
// terminated immediately with either an Error or Complete notification send to
// the subscribing client
//
// Note that this implementation is blocking. When no subcribers are present
// then the data can flow freely. But when there are subscribers, the observable
// goroutine is blocked until all subscribers have processed the next, error or
// complete notification.
func NewSubjectInt() SubjectInt {
	ch := multicast.NewChan(1, 16)

	observable := Observable(func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		ep, err := ch.NewEndpoint(0)
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
				return observer.Subscribed()
			}
			ep.Range(receive, 0)
		})
		observable(observe, subscribeOn, subscriber.Add(ep.Cancel))
	})

	observer := func(next int, err error, done bool) {
		if !ch.Closed() {
			if !done {
				ch.FastSend(next)
			} else {
				ch.Close(err)
			}
		}
	}

	return SubjectInt{observable.AsObservableInt(), observer}
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name ObservableIntSubscribeOn

// SubscribeOn specifies the scheduler an ObservableInt should use when it is
// subscribed to.
func (o ObservableInt) SubscribeOn(subscribeOn Scheduler) ObservableInt {
	observable := func(observe IntObserveFunc, _ Scheduler, subscriber Subscriber) {
		subscriber.OnWait(subscribeOn.Wait)
		o(observe, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Schedulers

func TrampolineScheduler() Scheduler	{ return scheduler.Trampoline }

func GoroutineScheduler() Scheduler	{ return scheduler.Goroutine }

//jig:name ObservableIntWait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Subscription is performed on the Trampoline scheduler.
func (o ObservableInt) Wait() (err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next int, e error, done bool) {
		if done {
			err = e
			subscriber.Unsubscribe()
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
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

//jig:name ObserveFuncMethods

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
	// Subscribed returns true when the subscription is currently valid.
	Subscribed() bool
}

//jig:name Create

// Create creates an Observable from scratch by calling observer methods
// programmatically.
func Create(f func(Observer)) Observable {
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			observer := func(next interface{}, err error, done bool) {
				if subscriber.Subscribed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				ObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

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
	options := &subscribeOptions{scheduler: TrampolineScheduler()}
	for _, setter := range setters {
		setter(options)
	}
	if options.subscriber == nil {
		options.subscriber = NewSubscriber()
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
// Subscribe by default is performed on the Trampoline scheduler.
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
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber
}
