package rx

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:template Connectable
//jig:needs Scheduler, Subscriber, Subscription

// Connectable provides the Connect method for a Multicaster.
type Connectable func(Scheduler, Subscriber)

// Connect instructs a multicaster to subscribe to its source and begin
// multicasting items to its subscribers.
func (c Connectable) Connect(subscribers ...Subscriber) Subscription {
	subscribers = append(subscribers, subscriber.New())
	scheduler := scheduler.MakeTrampoline()
	subscribers[0].OnWait(scheduler.Wait)
	c(scheduler, subscribers[0])
	return subscribers[0]
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
		terminated
	)
	var subject struct {
		state int32
		atomic.Value
	}
	subject.Store(factory())
	const (
		unsubscribed int32 = iota
		subscribed
	)
	var source struct {
		sync.Mutex
		state      int32
		subscriber Subscriber
	}
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subject.state, active, notifying) {
			if s, ok := subject.Load().(SubjectFoo); ok {
				s.FooObserver(next, err, done)
			}
			if !done {
				atomic.CompareAndSwapInt32(&subject.state, notifying, active)
			} else {
				if atomic.CompareAndSwapInt32(&subject.state, notifying, terminated) {
					source.Lock()
					source.subscriber.Unsubscribe()
					source.Unlock()
				}
			}
		}
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if s, ok := subject.Load().(SubjectFoo); ok {
			s.ObservableFoo(observe, subscribeOn, subscriber)
		}
	}
	connectable := func(subscribeOn Scheduler, subscriber Subscriber) {
		source.Lock()
		if atomic.CompareAndSwapInt32(&subject.state, terminated, active) {
			subject.Store(factory())
		}
		if atomic.CompareAndSwapInt32(&source.state, unsubscribed, subscribed) {
			source.subscriber = subscriber
			source.Unlock()
			o(observer, subscribeOn, subscriber)
			source.Lock()
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

// Publish uses the Multicast operator to control the subscription of a
// Subject to a source observable and turns the subject it into a connnectable
// observable. A Subject emits to an observer only those items that are emitted
// by the source Observable subsequent to the time of the observer subscribes.
//
// If the source completed and as a result the internal Subject terminated, then
// calling Connect again will replace the old Subject with a newly created one.
// So this Publish operator is re-connectable, unlike the RxJS 5 behavior that
// isn't. To simulate the RxJS 5 behavior use Publish().AutoConnect(1) this will
// connect on the first subscription but will never re-connect.
func (o ObservableFoo) Publish() FooMulticaster {
	return o.Multicast(NewSubjectFoo)
}

//jig:template Observable<Foo> PublishReplay
//jig:needs Observable<Foo> Multicast, NewReplaySubject<Foo>, <Foo>Multicaster

// Replay uses the Multicast operator to control the subscription of a
// ReplaySubject to a source observable and turns the subject into a
// connectable observable. A ReplaySubject emits to any observer all of the
// items that were emitted by the source observable, regardless of when the
// observer subscribes.
//
// If the source completed and as a result the internal ReplaySubject
// terminated, then calling Connect again will replace the old ReplaySubject
// with a newly created one.
func (o ObservableFoo) PublishReplay(bufferCapacity int, windowDuration time.Duration) FooMulticaster {
	factory := func() SubjectFoo {
		return NewReplaySubjectFoo(bufferCapacity, windowDuration)
	}
	return o.Multicast(factory)
}

//jig:template ErrRefCount
//jig:needs RxError

const ErrRefCountNeedsConcurrentScheduler = RxError("needs concurrent scheduler")

//jig:template <Foo>Multicaster RefCount
//jig:needs Observable<Foo>, ErrRefCount

// RefCount makes a FooMulticaster behave like an ordinary ObservableFoo. On
// first Subscribe it will call Connect on its FooMulticaster and when its last
// subscriber is Unsubscribed it will cancel the source connection by calling
// Unsubscribe on the subscription returned by the call to Connect.
func (o FooMulticaster) RefCount() ObservableFoo {
	var source struct {
		refcount     int32
		subscription Subscription
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, withSubscriber Subscriber) {
		if !subscribeOn.IsConcurrent() {
			var zero foo
			observe(zero, ErrRefCountNeedsConcurrentScheduler, true)
			return
		}
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			s := subscriber.New()
			o.Connectable(subscribeOn, s)
			source.subscription = s
		}
		withSubscriber.OnUnsubscribe(func() {
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				source.subscription.Unsubscribe()
			}
		})
		o.ObservableFoo(observe, subscribeOn, withSubscriber)
	}
	return observable
}

//jig:template ErrAutoConnect
//jig:needs RxError

const ErrAutoConnectInvalidCount = RxError("invalid count")
const ErrAutoConnectNeedsConcurrentScheduler = RxError("needs concurrent scheduler")

//jig:template <Foo>Multicaster AutoConnect
//jig:needs Observable<Foo>, Throw<Foo>, ErrAutoConnect

// AutoConnect makes a FooMulticaster behave like an ordinary ObservableFoo
// that automatically connects when the specified number of clients have
// subscribed to it. AutoConnect values should be larger or equal to 1.
// AutoConnect will throw an ErrInvalidCount if the count is out of range.
func (o FooMulticaster) AutoConnect(count int) ObservableFoo {
	if count < 1 {
		return ThrowFoo(ErrAutoConnectInvalidCount)
	}
	var refcount int32
	observable := func(observe FooObserver, subscribeOn Scheduler, withSubscriber Subscriber) {
		if !subscribeOn.IsConcurrent() {
			var zero foo
			observe(zero, ErrAutoConnectNeedsConcurrentScheduler, true)
			return
		}
		if atomic.AddInt32(&refcount, 1) == int32(count) {
			o.Connectable(subscribeOn, subscriber.New())
		}
		o.ObservableFoo(observe, subscribeOn, withSubscriber)
	}
	return observable
}
