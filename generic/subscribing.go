package rx

import (
	"fmt"

	"github.com/reactivego/subscriber"
)

//jig:template Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

//jig:template Observable<Foo> Print
//jig:needs Schedulers, Subscriber

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
func (o ObservableFoo) Println() (err error) {
	subscriber := subscriber.New()
	scheduler := CurrentGoroutineScheduler()
	observer := func(next foo, e error, done bool) {
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

//jig:template SubscribeOption
//jig:needs Schedulers, Subscriber

// SubscribeOption is an option that can be passed to the Subscribe method.
type SubscribeOption func(options *subscribeOptions)

type subscribeOptions struct {
	scheduler Scheduler
	subscriber Subscriber
	onSubscribe func(subscription Subscription)
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

//jig:template Observable<Foo> Subscribe
//jig:needs SubscribeOption

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
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
	return subscriber
}

//jig:template Observable<Foo> SubscribeNext
//jig:needs Observable<Foo> Subscribe

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscription.
func (o ObservableFoo) SubscribeNext(f func(next foo), options ...SubscribeOption) Subscription {
	return o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			f(next)
		}
	}, options...)
}

//jig:template Observable ToChan
//jig:needs Observable Subscribe

// ToChan returns a channel that emits interface{} values. If the source
// observable does not emit values but emits an error or complete, then the
// returned channel will enit any error and then close without emitting any
// values.
//
// Because the channel is fed by subscribing to the observable, ToChan would
// block when subscribed on the standard Trampoline scheduler which is initially
// synchronous. That's why the subscribing is done on the Goroutine scheduler.
//
// To cancel the subscription created internally by ToChan you will need access
// to the subscription used internnally by ToChan. To get at this subscription,
// pass the result of a call to option OnSubscribe(func(Subscription)) as a
// parameter to ToChan. On suscription the callback will be called with the
// subscription that was created.
func (o Observable) ToChan(options ...SubscribeOption) <-chan interface{} {
	scheduler := NewGoroutineScheduler()
	nextch := make(chan interface{}, 1)
	o.Subscribe(func(next interface{}, err error, done bool) {
		if !done {
			nextch <- next
		} else {
			if err != nil {
				nextch <- err
			}
			close(nextch)
		}
	}, SubscribeOn(scheduler, options...))
	return nextch
}

//jig:template Observable<Foo> ToChan
//jig:needs Observable<Foo> Subscribe
//jig:required-vars Foo

// ToChan returns a channel that emits foo values. If the source observable does
// not emit values but emits an error or complete, then the returned channel
// will close without emitting any values.
//
// There is no way to determine whether the observable feeding into the
// channel terminated with an error or completed normally.
// Because the channel is fed by subscribing to the observable, ToChan would
// block when subscribed on the standard Trampoline scheduler which is initially
// synchronous. That's why the subscribing is done on the Goroutine scheduler.
// It is not possible to cancel the subscription created internally by ToChan.
func (o ObservableFoo) ToChan(options ...SubscribeOption) <-chan foo {
	scheduler := NewGoroutineScheduler()
	nextch := make(chan foo, 1)
	o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			nextch <- next
		} else {
			close(nextch)
		}
	}, SubscribeOn(scheduler, options...))
	return nextch
}

//jig:template Observable<Foo> ToSingle
//jig:needs Observable<Foo> Subscribe

// ToSingle blocks until the ObservableFoo emits exactly one value or an error.
// The value and any error are returned.
//
// This function subscribes to the source observable on the Goroutine scheduler.
// The Goroutine scheduler works in more situations for complex chains of
// observables, like when merging the output of multiple observables.
func (o ObservableFoo) ToSingle(options ...SubscribeOption) (v foo, e error) {
	scheduler := NewGoroutineScheduler()
	o.Single().Subscribe(func(next foo, err error, done bool) {
		if !done {
			v = next
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, options...)).Wait()
	return v, e
}

//jig:template Observable<Foo> ToSlice
//jig:needs Observable<Foo> Subscribe

// ToSlice collects all values from the ObservableFoo into an slice. The
// complete slice and any error are returned.
//
// This function subscribes to the source observable on the Goroutine scheduler.
// The Goroutine scheduler works in more situations for complex chains of
// observables, like when merging the output of multiple observables.
func (o ObservableFoo) ToSlice(options ...SubscribeOption) (a []foo, e error) {
	scheduler := NewGoroutineScheduler()
	o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			a = append(a, next)
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, options...)).Wait()
	return a, e
}

//jig:template Observable<Foo> Wait
//jig:needs Observable<Foo> Subscribe

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
func (o ObservableFoo) Wait(options ...SubscribeOption) (e error) {
	o.Subscribe(func(next foo, err error, done bool) {
		if done {
			e = err
		}
	}, options...).Wait()
	return e
}
