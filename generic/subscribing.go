package rx

import (
	"fmt"
	"sync/atomic"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:template Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:template Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

//jig:template NewSubscriber
//jig:needs Subscriber

// NewSubscriber creates a new subscriber.
func NewSubscriber() Subscriber {
	return subscriber.New()
}

//jig:template Observable<Foo> AutoUnsubscribe<Foo>

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o ObservableFoo) AutoUnsubscribe() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next foo, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Println
//jig:needs Scheduler, Subscriber

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, err error, done bool) {
		if !done {
			fmt.Println(append(a, next)...)
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:template Observable<Foo> Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Wait() error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, err error, done bool) {
		if done {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:template Observable<Foo> ToSlice

// ToSlice collects all values from the ObservableFoo into an slice. The
// complete slice and any error are returned.
// ToSlice uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) ToSlice() (slice []foo, err error) {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, err error, done bool) {
		if !done {
			slice = append(slice, next)
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	err = subscriber.Wait()
	return
}

//jig:template Observable<Foo> ToSingle

// ToSingle blocks until the ObservableFoo emits exactly one value or an error.
// The value and any error are returned.
// ToSingle uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) ToSingle() (entry foo, err error) {
	o = o.Single()
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, err error, done bool) {
		if !done {
			entry = next
		} else {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	err = subscriber.Wait()
	return
}

//jig:template Observable ToChan

// ToChan returns a channel that emits interface{} values. If the source
// observable does not emit values but emits an error or complete, then the
// returned channel will enit any error and then close without emitting any
// values.
//
// ToChan uses the public scheduler.Goroutine variable for scheduling, because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To be able to cancel ToChan, you will need to
// create a subscriber yourself and pass it to ToChan as an argument.
func (o Observable) ToChan(subscribers ...Subscriber) <-chan interface{} {
	subscribers = append(subscribers, subscriber.New())
	scheduler := scheduler.Goroutine
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
				subscribers[0].Done(err)
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
//jig:required-vars Foo

// ToChan returns a channel that emits foo values. If the source observable does
// not emit values but emits an error or complete, then the returned channel
// will close without emitting any values.
//
// ToChan uses the public scheduler.Goroutine variable for scheduling, because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To be able to cancel ToChan, you will need to
// create a subscriber yourself and pass it to ToChan as an argument.
func (o ObservableFoo) ToChan(subscribers ...Subscriber) <-chan foo {
	subscribers = append(subscribers, subscriber.New())
	scheduler := scheduler.Goroutine
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
				subscribers[0].Done(err)
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

//jig:template Observable<Foo> Subscribe
//jig:needs Subscription

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Subscribe(observe FooObserver, schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero foo
			observe(zero, err, true)
			subscriber.Done(err)
		}
	}
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	o(observer, schedulers[0], subscriber)
	return subscriber
}

//jig:template Println<Foo>
//jig:needs <Foo>Observer

func PrintlnFoo(a ...interface{}) FooObserver {
	observer := func(next foo, err error, done bool) {
		if !done {
			fmt.Println(append(a, next)...)
		}
	}
	return observer
}
