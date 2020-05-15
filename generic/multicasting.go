package rx

import (
	"sync/atomic"
	"time"

	"github.com/reactivego/subscriber"
)

//jig:template Connectable
//jig:needs RxError, Subscription

// ErrDisconnect is sent to all subscribed observers when the subscription
// returned by Connect is cancelled by calling its Unsubscribe method.
const ErrDisconnect = RxError("disconnect")

// Connectable provides the Connect method for a Multicaster.
type Connectable func() Subscription

// Connect instructs a multicaster to subscribe to its source and begin
// multicasting items to its subscribers.
func (f Connectable) Connect() Subscription {
	return f()
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
	var subjectValue struct {
		state int32
		atomic.Value
	}
	subjectValue.Store(factory())
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subjectValue.state, active, notifying) {
			if s, ok := subjectValue.Load().(SubjectFoo); ok {
				s.FooObserver(next, err, done)
			}
			if !done {
				atomic.CompareAndSwapInt32(&subjectValue.state, notifying, active)
			} else {
				atomic.CompareAndSwapInt32(&subjectValue.state, notifying, terminated)
			}
		}
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if s, ok := subjectValue.Load().(SubjectFoo); ok {
			s.ObservableFoo(observe, subscribeOn, subscriber)
		}
	}
	const (
		unsubscribed int32 = iota
		subscribed
	)
	var subscriberValue struct {
		state int32
		atomic.Value
	}
	connectable := func() Subscription {
		if atomic.CompareAndSwapInt32(&subjectValue.state, terminated, active) {
			subjectValue.Store(factory())
		}
		if atomic.CompareAndSwapInt32(&subscriberValue.state, unsubscribed, subscribed) {
			subscriber := subscriber.New()
			o.Subscribe(observer, subscriber)
			subscriberValue.Store(subscriber)
			subscriber.OnUnsubscribe(func() {
				atomic.CompareAndSwapInt32(&subscriberValue.state, subscribed, unsubscribed)
				var zero foo
				observer(zero, ErrDisconnect, true)
			})
		}
		subscriber := subscriberValue.Load().(Subscriber)
		subscription := subscriber.Add(subscriber.Unsubscribe)
		subscription.OnWait(subscriber.Wait)
		return subscription
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

//jig:template <Foo>Multicaster RefCount
//jig:needs Observable<Foo>

// RefCount makes a FooMulticaster behave like an ordinary ObservableFoo. On
// first Subscribe it will call Connect on its FooMulticaster and when its last
// subscriber is Unsubscribed it will cancel the connection by calling
// Unsubscribe on the subscription returned by the call to Connect.
func (o FooMulticaster) RefCount() ObservableFoo {
	var (
		refcount   int32
		connection Subscription
	)
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&refcount, 1) == 1 {
			connection = o.Connect()
		}
		subscriber.OnUnsubscribe(func() {
			if atomic.AddInt32(&refcount, -1) == 0 {
				connection.Unsubscribe()
			}
		})
		o.ObservableFoo(observe, subscribeOn, subscriber)
	}
	return observable
}

//jig:template <Foo>Multicaster AutoConnect
//jig:needs Observable<Foo>

// AutoConnect makes a FooMulticaster behave like an ordinary ObservableFoo that
// automatically connects when the specified number of clients have subscribed
// to it. If count is 0, then AutoConnect will immediately call connect on the
// FooMulticaster before returning the ObservableFoo part of it.
func (o FooMulticaster) AutoConnect(count int) ObservableFoo {
	if count == 0 {
		o.Connect()
		return o.ObservableFoo
	}
	var refcount int32
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&refcount, 1) == int32(count) {
			o.Connect()
		}
		o.ObservableFoo(observe, subscribeOn, subscriber)
	}
	return observable
}
