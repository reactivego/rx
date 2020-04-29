package rx

import (
	"fmt"
	"sync/atomic"

	"github.com/reactivego/subscriber"
)

//jig:template Subscriber

// Subscriber is an alias for the subscriber.Subscriber interface type.
type Subscriber = subscriber.Subscriber

// NewSubscriber creates a new subscriber.
func NewSubscriber() Subscriber {
	return subscriber.New()
}

//jig:template Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

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
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
	return
}

//jig:template Observable<Foo> Wait
//jig:needs Schedulers, Subscriber

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Subscribing is performed on the Trampoline scheduler.
func (o ObservableFoo) Wait() (err error) {
	subscriber := NewSubscriber()
	scheduler := TrampolineScheduler()
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
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
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
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	subscriber.Wait()
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

//jig:template Observable<Foo> Subscribe
//jig:needs Schedulers, Subscriber, Subscription

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe by default is performed on the Trampoline scheduler.
func (o ObservableFoo) Subscribe(observe FooObserveFunc, subscribers ...Subscriber) Subscription {
	subscribers = append(subscribers, NewSubscriber())
	scheduler := TrampolineScheduler()
	observer := func(next foo, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zeroFoo, err, true)
			subscribers[0].Unsubscribe()
		}
	}
	subscribers[0].OnWait(scheduler.Wait)
	o(observer, scheduler, subscribers[0])
	return subscribers[0]
}
