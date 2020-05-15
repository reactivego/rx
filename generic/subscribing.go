package rx

import (
	"fmt"
	"sync/atomic"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
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

//jig:template Observable<Foo> Println
//jig:needs Scheduler, Subscriber

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Println(a ...interface{}) (err error) {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, e error, done bool) {
		if !done {
			fmt.Println(append(a, next)...)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
}

//jig:template Observable<Foo> Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Wait() (err error) {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, e error, done bool) {
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

//jig:template Observable<Foo> ToSlice

// ToSlice collects all values from the ObservableFoo into an slice. The
// complete slice and any error are returned.
// ToSlice uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) ToSlice() (slice []foo, err error) {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, e error, done bool) {
		if !done {
			slice = append(slice, next)
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
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
	observer := func(next foo, e error, done bool) {
		if !done {
			entry = next
		} else {
			err = e
			subscriber.Unsubscribe()
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
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
// by the calling code directly. To cancel ToChan you will need to supply a
// subscriber that you hold on to.
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
//jig:required-vars Foo

// ToChan returns a channel that emits foo values. If the source observable does
// not emit values but emits an error or complete, then the returned channel
// will close without emitting any values.
//
// ToChan uses the public scheduler.Goroutine variable for scheduling, because
// it needs the concurrency so the returned channel can be used by used
// by the calling code directly. To cancel ToChan you will need to supply a
// subscriber that you hold on to.
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

//jig:template Observable<Foo> Subscribe
//jig:needs Subscription

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableFoo) Subscribe(observe FooObserver, subscribers ...Subscriber) Subscription {
	subscribers = append(subscribers, subscriber.New())
	scheduler := scheduler.MakeTrampoline()
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zeroFoo foo
			observe(zeroFoo, err, true)
			subscribers[0].Unsubscribe()
		}
	}
	subscribers[0].OnWait(scheduler.Wait)
	o(observer, scheduler, subscribers[0])
	return subscribers[0]
}
