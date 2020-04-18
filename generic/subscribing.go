package rx

import (
	"fmt"
	"sync/atomic"

	"github.com/reactivego/subscriber"
)

//jig:template Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

// NewSubscriber creates a new subscriber.
func NewSubscriber() Subscriber {
	return subscriber.New()
}

//jig:template Observable<Foo> Println
//jig:needs Schedulers, Subscriber

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println is performed on the Trampoline scheduler.
func (o ObservableFoo) Println() (err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next foo, e error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	scheduler.Wait()
	return
}

//jig:template Observable<Foo> Wait
//jig:needs Schedulers, Subscriber

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Subscription is performed on the Trampoline scheduler.
func (o ObservableFoo) Wait() (err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next foo, e error, done bool) {
		if done {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	scheduler.Wait()
	return
}

//jig:template Observable<Foo> ToSlice
//jig:needs Schedulers, Subscriber

// ToSlice collects all values from the ObservableFoo into an slice. The
// complete slice and any error are returned.
func (o ObservableFoo) ToSlice() (slice []foo, err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next foo, e error, done bool) {
		if !done {
			slice = append(slice, next)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	scheduler.Wait()
	return
}

//jig:template Observable<Foo> ToSingle
//jig:needs Schedulers, Subscriber

// ToSingle blocks until the ObservableFoo emits exactly one value or an error.
// The value and any error are returned.
func (o ObservableFoo) ToSingle() (entry foo, err error) {
	o = o.Single()
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
	observer := func(next foo, e error, done bool) {
		if !done {
			entry = next
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	scheduler.Wait()
	return
}

//jig:template Observable ToChan
//jig:needs Schedulers, Subscriber

// ToChan returns a channel that emits interface{} values. If the source
// observable does not emit values but emits an error or complete, then the
// returned channel will enit any error and then close without emitting any
// values.
//
// This method subscribes to the observable on the Goroutine scheduler because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To cancel ToChan you will need to supply a
// subscriber that you hold on to.
func (o Observable) ToChan(subscribers ...Subscriber) <-chan interface{} {
	scheduler := GoroutineScheduler()
	subscribers = append(subscribers, NewSubscriber())
	donech := make(chan struct{})
	nextch := make(chan interface{})
	const (
		idle = iota
		busy
		closed
	)
	state := int32(idle)
	observer := func(next interface{}, err error, done bool) {
		if atomic.CompareAndSwapInt32(&state, idle, busy) {
			if err != nil {
				next = err
			}
			if !done || err != nil {
				select {
				case <-donech:
					atomic.StoreInt32(&state, closed)
				default:
					select {
					case <-donech:
						atomic.StoreInt32(&state, closed)
					case nextch <- next:
					}
				}
			}
			if done {
				atomic.StoreInt32(&state, closed)
				subscribers[0].Unsubscribe()
			}
			if !atomic.CompareAndSwapInt32(&state, busy, idle) {
				close(nextch)
			}
		}
	}
	subscribers[0].OnUnsubscribe(func() {
		close(donech)
		if atomic.CompareAndSwapInt32(&state, busy, closed) {
			return
		}
		if atomic.CompareAndSwapInt32(&state, idle, closed) {
			close(nextch)
			return
		}
	})
	o(observer, scheduler, subscribers[0])
	return nextch
}

//jig:template Observable<Foo> ToChan
//jig:needs Schedulers, Subscriber
//jig:required-vars Foo

// ToChan returns a channel that emits foo values. If the source observable does
// not emit values but emits an error or complete, then the returned channel
// will close without emitting any values.
//
// This method subscribes to the observable on the Goroutine scheduler because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To cancel ToChan you will need to supply a
// subscriber that you hold on to.
func (o ObservableFoo) ToChan(subscribers ...Subscriber) <-chan foo {
	scheduler := GoroutineScheduler()
	subscribers = append(subscribers, NewSubscriber())
	donech := make(chan struct{})
	nextch := make(chan foo)
	const (
		idle = iota
		busy
		closed
	)
	state := int32(idle)
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&state, idle, busy) {
			if !done {
				select {
				case <-donech:
					atomic.StoreInt32(&state, closed)
				default:
					select {
					case <-donech:
						atomic.StoreInt32(&state, closed)
					case nextch <- next:
					}
				}
			} else {
				atomic.StoreInt32(&state, closed)
				subscribers[0].Unsubscribe()
			}
			if !atomic.CompareAndSwapInt32(&state, busy, idle) {
				close(nextch)
			}
		}
	}
	subscribers[0].OnUnsubscribe(func() {
		close(donech)
		if atomic.CompareAndSwapInt32(&state, busy, closed) {
			return
		}
		if atomic.CompareAndSwapInt32(&state, idle, closed) {
			close(nextch)
			return
		}
	})
	o(observer, scheduler, subscribers[0])
	return nextch
}

//jig:template SubscribeOption
//jig:needs Schedulers, Subscriber

// SubscribeOption is an option that can be passed to the Subscribe method.
type SubscribeOption func(options *subscribeOptions)

type subscribeOptions struct {
	scheduler     Scheduler
	subscriber    Subscriber
	onSubscribe   func(subscription Subscription)
	onUnsubscribe func()
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

//jig:template Observable<Foo> Subscribe
//jig:needs SubscribeOption

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe by default is performed on the Trampoline scheduler.
func (o ObservableFoo) Subscribe(observe FooObserveFunc, options ...SubscribeOption) Subscription {
	scheduler, subscriber := newSchedulerAndSubscriber(options)
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFoo, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, scheduler, subscriber)
	subscriber.OnWait(scheduler.Wait)
	return subscriber
}

//jig:template Observable<Foo> SubscribeNext
//jig:needs Observable<Foo> Subscribe

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscription.
// SubscribeNext by default is performed on the Trampoline scheduler.
func (o ObservableFoo) SubscribeNext(f func(next foo), options ...SubscribeOption) Subscription {
	return o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			f(next)
		}
	}, options...)
}
