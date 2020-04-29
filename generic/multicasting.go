package rx

import (
	"sync/atomic"
	"time"

	"github.com/reactivego/subscriber"
)

//jig:template Connectable<Foo>
//jig:embeds Observable<Foo>

// ConnectableFoo is an ObservableFoo that has an additional method Connect()
// used to Subscribe to the parent observable and then multicasting values to
// all subscribers of ConnectableFoo.
type ConnectableFoo struct {
	ObservableFoo
	connect func() Subscription
}

//jig:template Connectable<Foo> Connect
//jig:needs Connectable<Foo>

// Connect instructs a connectable Observable to begin emitting items to its
// subscribers. All values will then be passed on to the observers that
// subscribed to this connectable observable
func (c ConnectableFoo) Connect() Subscription {
	return c.connect()
}

//jig:template Observable<Foo> Multicast
//jig:needs Observable<Foo> Subscribe, Connectable<Foo>, Schedulers

// Multicast converts an ordinary Observable into a connectable Observable.
// A connectable observable will only start emitting values after its Connect
// method has been called. The factory method passed in should return a
// new SubjectFoo that implements the actual multicasting behavior.
func (o ObservableFoo) Multicast(factory func() SubjectFoo) ConnectableFoo {
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
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if s, ok := subjectValue.Load().(SubjectFoo); ok {
			s.ObservableFoo(observe, subscribeOn, subscriber)
		}
	}
	observer := func(next foo, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subjectValue.state, active, notifying) {
			if s, ok := subjectValue.Load().(SubjectFoo); ok {
				s.FooObserveFunc(next, err, done)
			}
			if !done {
				atomic.CompareAndSwapInt32(&subjectValue.state, notifying, active)
			} else {
				atomic.CompareAndSwapInt32(&subjectValue.state, notifying, terminated)
			}
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
	connect := func() Subscription {
		if atomic.CompareAndSwapInt32(&subjectValue.state, terminated, active) {
			subjectValue.Store(factory())
		}
		if atomic.CompareAndSwapInt32(&subscriberValue.state, unsubscribed, subscribed) {
			scheduler := GoroutineScheduler()
			subscriber := subscriber.New()
			o.SubscribeOn(scheduler).Subscribe(observer, subscriber)
			subscriberValue.Store(subscriber)
			subscriber.OnUnsubscribe(func() {
				atomic.CompareAndSwapInt32(&subscriberValue.state, subscribed, unsubscribed)
			})
		}
		subscription := subscriberValue.Load().(Subscriber)
		return subscription.Add(func() { subscription.Unsubscribe() })
	}
	return ConnectableFoo{ObservableFoo: observable, connect: connect}
}

//jig:template Observable<Foo> Publish
//jig:needs Observable<Foo> Multicast, NewSubject<Foo>, Connectable<Foo>

// Publish uses Multicast to control the subscription of a Subject to a
// source observable and turns the subject it into a connnectable observable.
// A Subject emits to an observer only those items that are emitted by
// the source Observable subsequent to the time of the subscription.
//
// If the source completed and as a result the internal Subject terminated, then
// calling Connect again will replace the old Subject with a newly created one.
// So this Publish operator is re-connectable, unlike the RxJS 5 behavior that
// isn't. To simulate the RxJS 5 behavior use Publish().AutoConnect(1) this will
// connect on the first subscription but will never re-connect.
func (o ObservableFoo) Publish() ConnectableFoo {
	return o.Multicast(NewSubjectFoo)
}

//jig:template Observable<Foo> PublishReplay
//jig:needs Observable<Foo> Multicast, NewReplaySubject<Foo>, Connectable<Foo>

// Replay uses Multicast to control the subscription of a ReplaySubject to a
// source observable and turns the subject into a connectable observable.
// A ReplaySubject emits to any observer all of the items that were emitted by
// the source observable, regardless of when the observer subscribes.
//
// If the source completed and as a result the internal ReplaySubject
// terminated, then calling Connect again will replace the old ReplaySubject
// with a newly created one.
func (o ObservableFoo) PublishReplay(bufferCapacity int, windowDuration time.Duration) ConnectableFoo {
	factory := func() SubjectFoo {
		return NewReplaySubjectFoo(bufferCapacity, windowDuration)
	}
	return o.Multicast(factory)
}

//jig:template Connectable<Foo> RefCount
//jig:needs Observable<Foo>

// RefCount makes a ConnectableFoo behave like an ordinary ObservableFoo. On
// first Subscribe it will call Connect on its ConnectableFoo and when its last
// subscriber is Unsubscribed it will cancel the connection by calling
// Unsubscribe on the subscription returned by the call to Connect.
func (o ConnectableFoo) RefCount() ObservableFoo {
	var (
		refcount   int32
		connection Subscription
	)
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&refcount, 1) == 1 {
			connection = o.connect()
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

//jig:template Connectable<Foo> AutoConnect
//jig:needs Observable<Foo>

// AutoConnect makes a ConnectableFoo behave like an ordinary ObservableFoo that
// automatically connects when the specified number of clients subscribe to it.
// If count is 0, then AutoConnect will immediately call connect on the
// ConnectableFoo before returning the ObservableFoo part of the ConnectableFoo.
func (o ConnectableFoo) AutoConnect(count int) ObservableFoo {
	if count == 0 {
		o.connect()
		return o.ObservableFoo
	}
	var refcount int32
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&refcount, 1) == int32(count) {
			o.connect()
		}
		o.ObservableFoo(observe, subscribeOn, subscriber)
	}
	return observable
}
