package rx

import (
	"fmt"

	"github.com/reactivego/subscriber"
)

//jig:template Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber subscriber.Subscriber

//jig:template SubscribeOptions
//jig:needs Scheduler, Subscriber

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription subscriber.Subscription

// SubscribeOptions is a struct with options for Subscribe related methods.
type SubscribeOptions struct {
	// SubscribeOn is the scheduler to run the observable subscription on.
	SubscribeOn Scheduler
	// OnSubscribe is called right after the subscription is created and before
	// subscribing continues further.
	OnSubscribe func(subscription Subscription)
	// OnUnsubscribe is called by the subscription to notify the client that the
	// subscription has been canceled.
	OnUnsubscribe func()
}

// NewSubscriber will return a newly created subscriber. Before returning the
// subscription the OnSubscribe callback (if set) will already have been called.
func (options SubscribeOptions) NewSubscriber() Subscriber {
	subscriber := subscriber.New()
	subscriber.OnUnsubscribe(options.OnUnsubscribe)
	if options.OnSubscribe != nil {
		options.OnSubscribe(subscriber)
	}
	return subscriber
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

//jig:template Observable<Foo> Subscribe
//jig:needs NewScheduler, SubscribeOptions

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableFoo) Subscribe(observe FooObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFoo, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:template Observable<Foo> SubscribeNext
//jig:needs Observable<Foo> Subscribe

// SubscribeNext operates upon the emissions from an Observable only.
// This method returns a Subscription.
func (o ObservableFoo) SubscribeNext(f func(next foo), setters ...SubscribeOptionSetter) Subscription {
	return o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			f(next)
		}
	}, setters...)
}

//jig:template Observable<Foo> Println
//jig:needs Observable<Foo> Subscribe

// Println subscribes to the Observable and prints every item to os.Stdout while
// it waits for completion or error. Returns either the error or nil when the
// Observable completed normally.
func (o ObservableFoo) Println(setters ...SubscribeOptionSetter) (e error) {
	o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			e = err
		}
	}, setters...).Wait()
	return e
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
func (o Observable) ToChan(setters ...SubscribeOptionSetter) <-chan interface{} {
	scheduler := NewGoroutine()
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
	}, SubscribeOn(scheduler, setters...))
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
func (o ObservableFoo) ToChan(setters ...SubscribeOptionSetter) <-chan foo {
	scheduler := NewGoroutine()
	nextch := make(chan foo, 1)
	o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			nextch <- next
		} else {
			close(nextch)
		}
	}, SubscribeOn(scheduler, setters...))
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
func (o ObservableFoo) ToSingle(setters ...SubscribeOptionSetter) (v foo, e error) {
	scheduler := NewGoroutine()
	o.Single().Subscribe(func(next foo, err error, done bool) {
		if !done {
			v = next
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, setters...)).Wait()
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
func (o ObservableFoo) ToSlice(setters ...SubscribeOptionSetter) (a []foo, e error) {
	scheduler := NewGoroutine()
	o.Subscribe(func(next foo, err error, done bool) {
		if !done {
			a = append(a, next)
		} else {
			e = err
		}
	}, SubscribeOn(scheduler, setters...)).Wait()
	return a, e
}

//jig:template Observable<Foo> Wait
//jig:needs Observable<Foo> Subscribe

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
func (o ObservableFoo) Wait(setters ...SubscribeOptionSetter) (e error) {
	o.Subscribe(func(next foo, err error, done bool) {
		if done {
			e = err
		}
	}, setters...).Wait()
	return e
}
