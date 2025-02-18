package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

//jig:template Connectable
//jig:needs Scheduler, Subscriber, Subscription

// Connectable provides the Connect method for a Multicaster.
type Connectable func(Scheduler, Subscriber)

// Connect instructs a multicaster to subscribe to its source and begin
// multicasting items to its subscribers. Connect accepts an optional
// scheduler argument.
func (c Connectable) Connect(schedulers ...Scheduler) Subscription {
	subscriber := NewSubscriber()
	schedulers = append(schedulers, NewScheduler())
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	c(schedulers[0], subscriber)
	return subscriber
}

//jig:template <Foo>Multicaster
//jig:embeds Observable<Foo>, Connectable

// FooMulticaster is a multicasting connectable observable. One or more
// FooObservers can subscribe to it simultaneously. It will subscribe to the
// source ObservableFoo when Connect is called. After that, every emission
// from the source is multcast to all subscribed FooObservers.
type FooMulticaster struct {
	ObservableFoo
	Connectable
}

//jig:template Observable<Foo> Multicast
//jig:needs Observable<Foo> Subscribe, <Foo>Multicaster

// Multicast converts an ordinary observable into a multicasting connectable
// observable or multicaster for short. A multicaster will only start emitting
// values after its Connect method has been called. The factory method passed
// in should return a new SubjectFoo that implements the actual multicasting
// behavior.
func (o ObservableFoo) Multicast(factory func() SubjectFoo) FooMulticaster {
	const (
		active int32 = iota
		notifying
		erred
		completed
	)
	var subject struct {
		state int32
		atomic.Value
		count int32
	}
	const (
		unsubscribed int32 = iota
		subscribed
	)
	var source struct {
		sync.Mutex
		state      int32
		subscriber Subscriber
	}
	subject.Store(factory())
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subject.state, active, notifying) {
			if s, ok := subject.Load().(SubjectFoo); ok {
				s.FooObserver(next, err, done)
			}
			switch {
			case !done:
				atomic.CompareAndSwapInt32(&subject.state, notifying, active)
			case err != nil:
				if atomic.CompareAndSwapInt32(&subject.state, notifying, erred) {
					source.subscriber.Done(err)
				}
			default:
				if atomic.CompareAndSwapInt32(&subject.state, notifying, completed) {
					source.subscriber.Done(nil)
				}
			}
		}
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&subject.count, 1) == 1 {
			if atomic.CompareAndSwapInt32(&subject.state, erred, active) {
				subject.Store(factory())
			}
		}
		if s, ok := subject.Load().(SubjectFoo); ok {
			s.ObservableFoo(observe, subscribeOn, subscriber)
		}
		subscriber.OnUnsubscribe(func() {
			atomic.AddInt32(&subject.count, -1)
		})
	}
	connectable := func(subscribeOn Scheduler, subscriber Subscriber) {
		source.Lock()
		if atomic.CompareAndSwapInt32(&source.state, unsubscribed, subscribed) {
			source.subscriber = subscriber
			o(observer, subscribeOn, subscriber)
			subscriber.OnUnsubscribe(func() {
				atomic.CompareAndSwapInt32(&source.state, subscribed, unsubscribed)
			})
		} else {
			source.subscriber.OnUnsubscribe(subscriber.Unsubscribe)
			subscriber.OnUnsubscribe(source.subscriber.Unsubscribe)
		}
		source.Unlock()
	}
	return FooMulticaster{ObservableFoo: observable, Connectable: connectable}
}

//jig:template Observable<Foo> Publish
//jig:needs Observable<Foo> Multicast, NewSubject<Foo>, <Foo>Multicaster

// Publish returns a FooMulticaster for a Subject to an underlying
// FooObservable and turns the subject into a connnectable observable. A
// Subject emits to an observer only those items that are emitted by the
// underlying FooObservable subsequent to the time of the observer subscribes.
// When the underlying FooObervable terminates with an error, then subscribed
// observers will receive that error. After all observers have unsubscribed
// due to an error, the FooMulticaster does an internal reset just before the
// next observer subscribes. So this Publish operator is re-connectable,
// unlike the RxJS 5 behavior that isn't. To simulate the RxJS 5 behavior use
// Publish().AutoConnect(1) this will connect on the first subscription but
// will never re-connect.
func (o ObservableFoo) Publish() FooMulticaster {
	return o.Multicast(NewSubjectFoo)
}

//jig:template Observable<Foo> PublishReplay
//jig:needs Observable<Foo> Multicast, NewReplaySubject<Foo>, <Foo>Multicaster

// PublishReplay returns a FooMulticaster for a ReplaySubject to an underlying
// FooObservable and turns the subject into a connectable observable. A
// ReplaySubject emits to any observer all of the items that were emitted by
// the source observable, regardless of when the observer subscribes. When the
// underlying FooObervable terminates with an error, then subscribed observers
// will receive that error. After all observers have unsubscribed due to an
// error, the FooMulticaster does an internal reset just before the next
// observer subscribes.
func (o ObservableFoo) PublishReplay(bufferCapacity int, windowDuration time.Duration) FooMulticaster {
	factory := func() SubjectFoo {
		return NewReplaySubjectFoo(bufferCapacity, windowDuration)
	}
	return o.Multicast(factory)
}

//jig:template Observable<Foo> PublishBehavior
//jig:needs Observable<Foo> Multicast, NewBehaviorSubject<Foo>, <Foo>Multicaster

// PublishBehavior returns a FooMulticaster that shares a single subscription
// to the underlying FooObservable returning an initial value or the last
// value emitted by the underlying FooObservable. When the underlying
// FooObervable terminates with an error, then subscribed observers will
// receive that error. After all observers have unsubscribed due to an error,
// the FooMulticaster does an internal reset just before the next observer
// subscribes.
func (o ObservableFoo) PublishBehavior(a foo) FooMulticaster {
	factory := func() SubjectFoo {
		return NewBehaviorSubjectFoo(a)
	}
	return o.Multicast(factory)
}

//jig:template Observable<Foo> PublishLast
//jig:needs Observable<Foo> Multicast, NewAsyncSubject<Foo>, <Foo>Multicaster

// PublishLast returns a FooMulticaster that shares a single subscription to
// the underlying FooObservable containing only the last value emitted before
// it completes. When the underlying FooObervable terminates with an error,
// then subscribed observers will receive only that error (and no value).
// After all observers have unsubscribed due to an error, the FooMulticaster
// does an internal reset just before the next observer subscribes.
func (o ObservableFoo) PublishLast() FooMulticaster {
	return o.Multicast(NewAsyncSubjectFoo)
}

//jig:template <Foo>Multicaster RefCount
//jig:needs Observable<Foo>

// RefCount makes a FooMulticaster behave like an ordinary ObservableFoo. On
// first Subscribe it will call Connect on its FooMulticaster and when its last
// subscriber is Unsubscribed it will cancel the source connection by calling
// Unsubscribe on the subscription returned by the call to Connect.
func (o FooMulticaster) RefCount() ObservableFoo {
	var source struct {
		sync.Mutex
		refcount   int32
		subscriber Subscriber
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, withSubscriber Subscriber) {
		withSubscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				source.subscriber.Unsubscribe()
			}
			source.Unlock()
		})
		o.ObservableFoo(observe, subscribeOn, withSubscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			source.subscriber = NewSubscriber()
			source.Unlock()
			o.Connectable(subscribeOn, source.subscriber)
			source.Lock()
		}
		source.Unlock()
	}
	return observable
}

//jig:template InvalidCount
//jig:needs RxError

const InvalidCount = RxError("invalid count")

//jig:template <Foo>Multicaster AutoConnect
//jig:needs Observable<Foo>, Throw<Foo>, InvalidCount

// AutoConnect makes a FooMulticaster behave like an ordinary ObservableFoo
// that automatically connects the multicaster to its source when the
// specified number of observers have subscribed to it. If the count is less
// than 1 it will return a ThrowFoo(InvalidCount). After connecting, when the
// number of subscribed observers eventually drops to 0, AutoConnect will
// cancel the source connection if it hasn't terminated yet. When subsequently
// the next observer subscribes, AutoConnect will connect to the source only
// when it was previously canceled or because the source terminated with an
// error. So it will not reconnect when the source completed succesfully. This
// specific behavior allows for implementing a caching observable that can be
// retried until it succeeds. Another thing to notice is that AutoConnect will
// disconnect an active connection when the number of observers drops to zero.
// The reason for this is that not doing so would leak a task and leave it
// hanging in the scheduler.
func (o FooMulticaster) AutoConnect(count int) ObservableFoo {
	if count < 1 {
		return ThrowFoo(InvalidCount)
	}
	var source struct {
		sync.Mutex
		refcount   int32
		subscriber Subscriber
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, withSubscriber Subscriber) {
		withSubscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				if source.subscriber != nil {
					source.subscriber.Unsubscribe()
				}
			}
			source.Unlock()
		})
		o.ObservableFoo(observe, subscribeOn, withSubscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == int32(count) {
			if source.subscriber == nil || source.subscriber.Error() != nil {
				source.subscriber = NewSubscriber()
				source.Unlock()
				o.Connectable(subscribeOn, source.subscriber)
				source.Lock()
			}
		}
		source.Unlock()
	}
	return observable
}
