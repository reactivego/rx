// Code generated by jig; DO NOT EDIT.

//go:generate jig

package Retry

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:name Error

// Error signals an error condition.
type Error func(error)

//jig:name Complete

// Complete signals that no more data is to be expected.
type Complete func()

//jig:name Canceled

// Canceled returns true when the observer has unsubscribed.
type Canceled func() bool

//jig:name NextInt

// NextInt can be called to emit the next value to the IntObserver.
type NextInt func(int)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

//jig:name IntObserver

// IntObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type IntObserver func(next int, err error, done bool)

//jig:name ObservableInt

// ObservableInt is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableInt func(IntObserver, Scheduler, Subscriber)

//jig:name CreateInt

// CreateInt provides a way of creating an ObservableInt from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateInt will be called once
// to implement the observable. It is provided with a NextInt, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func CreateInt(create func(NextInt, Error, Complete, Canceled)) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next int) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				if subscriber.Subscribed() {
					var zero int
					observe(zero, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					var zero int
					observe(zero, nil, true)
				}
			}
			x := func() bool {
				return !subscriber.Subscribed()
			}
			create(n, e, c, x)
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name GoroutineScheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

//jig:name MakeTrampolineScheduler

func MakeTrampolineScheduler() Scheduler {
	return scheduler.MakeTrampoline()
}

//jig:name Observer

// Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Observer func(next interface{}, err error, done bool)

//jig:name Observable

// Observable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type Observable func(Observer, Scheduler, Subscriber)

//jig:name Defer

// Defer does not create the Observable until the observer subscribes.
// It creates a fresh Observable for each subscribing observer. Use it to
// create observables that maintain separate state per subscription.
func Defer(factory func() Observable) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name Throw

// Throw creates an Observable that emits no items and terminates with an
// error.
func Throw(err error) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				var zero interface{}
				observe(zero, err, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name From

// From creates an Observable from multiple interface{} values passed in.
func From(slice ...interface{}) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					var zero interface{}
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableInt_Retry

// Retry if a source ObservableInt sends an error notification, resubscribe to
// it in the hopes that it will complete without error. If count is zero or
// negative, the retry count will be effectively infinite. The scheduler
// passed when subscribing is used by Retry to schedule any retry attempt. The
// time between retries is 1 millisecond, so retry frequency is 1 kHz. Any
// SubscribeOn operators should be called after Retry to prevent lockups
// caused by mixing different schedulers in the same subscription for retrying
// and subscribing.
func (o ObservableInt) Retry(count ...int) ObservableInt {
	count = append(count, int(^uint(0)>>1))
	if count[0] <= 0 {
		count = []int{int(^uint(0) >> 1)}
	}
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var retry struct {
			count		int
			observer	IntObserver
			subscriber	Subscriber
			resubscribe	func()
		}
		retry.count = count[0]
		retry.observer = func(next int, err error, done bool) {
			if err != nil && retry.count > 0 {
				retry.count--
				retry.subscriber.Done(err)
				subscribeOn.ScheduleFuture(1*time.Millisecond, retry.resubscribe)
			} else {
				observe(next, err, done)
			}
		}
		retry.resubscribe = func() {
			if subscriber.Subscribed() {
				retry.subscriber = subscriber.Add()
				if !subscribeOn.IsConcurrent() {
					retry.subscriber.OnWait(subscribeOn.Wait)
					subscriber.OnWait(func() { retry.subscriber.Wait() })
				}
				o(retry.observer, subscribeOn, retry.subscriber)
			}
		}
		retry.resubscribe()
	}
	return observable
}

//jig:name Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:name Observable_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o Observable) Subscribe(observe Observer, schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero interface{}
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

//jig:name Connectable

// Connectable provides the Connect method for a Multicaster.
type Connectable func(Scheduler, Subscriber)

// Connect instructs a multicaster to subscribe to its source and begin
// multicasting items to its subscribers. Connect accepts an optional
// scheduler argument.
func (c Connectable) Connect(schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	c(schedulers[0], subscriber)
	return subscriber
}

//jig:name Multicaster

// Multicaster is a multicasting connectable observable. One or more
// Observers can subscribe to it simultaneously. It will subscribe to the
// source Observable when Connect is called. After that, every emission
// from the source is multcast to all subscribed Observers.
type Multicaster struct {
	Observable
	Connectable
}

//jig:name Observable_Multicast

// Multicast converts an ordinary observable into a multicasting connectable
// observable or multicaster for short. A multicaster will only start emitting
// values after its Connect method has been called. The factory method passed
// in should return a new Subject that implements the actual multicasting
// behavior.
func (o Observable) Multicast(factory func() Subject) Multicaster {
	const (
		active	int32	= iota
		notifying
		erred
		completed
	)
	var subject struct {
		state	int32
		atomic.Value
		count	int32
	}
	const (
		unsubscribed	int32	= iota
		subscribed
	)
	var source struct {
		sync.Mutex
		state		int32
		subscriber	Subscriber
	}
	subject.Store(factory())
	observer := func(next interface{}, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subject.state, active, notifying) {
			if s, ok := subject.Load().(Subject); ok {
				s.Observer(next, err, done)
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
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&subject.count, 1) == 1 {
			if atomic.CompareAndSwapInt32(&subject.state, erred, active) {
				subject.Store(factory())
			}
		}
		if s, ok := subject.Load().(Subject); ok {
			s.Observable(observe, subscribeOn, subscriber)
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
	return Multicaster{Observable: observable, Connectable: connectable}
}

//jig:name MakeObserverObservable

const OutOfSubscriptions = RxError("out of subscriptions")

// MakeObserverObservable turns an observer into a multicasting and buffering
// observable. Both the observer and the obeservable are returned. These are
// then used as the core of any Subject implementation. The Observer side is
// used to pass items into the buffering multicaster. This then multicasts the
// items to every Observer that subscribes to the returned Observable.
//
//	age     age below which items are kept to replay to a new subscriber.
//	length  length of the item buffer, number of items kept to replay to a new subscriber.
//	[cap]   Capacity of the item buffer, number of items that can be observed before blocking.
//	[scap]  Capacity of the subscription list, max number of simultaneous subscribers.
func MakeObserverObservable(age time.Duration, length int, capacity ...int) (Observer, Observable) {
	const (
		ms	= time.Millisecond
		us	= time.Microsecond
	)

	type subscription struct {
		cursor		uint64
		state		uint64		// active, canceled, closed
		activated	time.Time	// track activity to deterime backoff
		subscribeOn	Scheduler
	}

	// cursor
	const (
		maxuint64 uint64 = math.MaxUint64	// park unused cursor at maxuint64
	)

	// state
	const (
		active	uint64	= iota
		canceled
		closing
		closed
	)

	type subscriptions struct {
		sync.Mutex
		*sync.Cond
		entries	[]subscription
		access	uint32	// unlocked, locked
	}

	// access
	const (
		unlocked	uint32	= iota
		locked
	)

	type item struct {
		Value	interface{}
		At	time.Time
	}

	type buffer struct {
		age	time.Duration
		keep	uint64
		mod	uint64
		size	uint64

		items	[]item
		begin	uint64
		end	uint64
		commit	uint64
		state	uint64	// active, closed

		subscriptions	subscriptions

		err	error
	}

	make := func(age time.Duration, length int, capacity ...int) *buffer {
		if length < 0 {
			length = 0
		}
		keep := uint64(length)

		cap, scap := length, 32
		switch {
		case len(capacity) >= 2:
			cap, scap = capacity[0], capacity[1]
		case len(capacity) == 1:
			cap = capacity[0]
		}
		if cap < length {
			cap = length
		}
		if cap == 0 {
			cap = 1
		}
		size := uint64(1) << uint(math.Ceil(math.Log2(float64(cap))))

		if scap < 1 {
			scap = 1
		}
		buf := &buffer{
			age:	age,
			keep:	keep,
			mod:	size - 1,
			size:	size,

			items:	make([]item, size),
			end:	size,
			subscriptions: subscriptions{
				entries: make([]subscription, 0, scap),
			},
		}
		buf.subscriptions.Cond = sync.NewCond(&buf.subscriptions.Mutex)
		return buf
	}
	buf := make(age, length, capacity...)

	accessSubscriptions := func(access func([]subscription)) bool {
		gosched := false
		for !atomic.CompareAndSwapUint32(&buf.subscriptions.access, unlocked, locked) {
			runtime.Gosched()
			gosched = true
		}
		access(buf.subscriptions.entries)
		atomic.StoreUint32(&buf.subscriptions.access, unlocked)
		return gosched
	}

	send := func(value interface{}) {
		for buf.commit == buf.end {
			full := false
			subscribeOn := Scheduler(nil)
			gosched := accessSubscriptions(func(subscriptions []subscription) {
				slowest := maxuint64
				for i := range subscriptions {
					current := atomic.LoadUint64(&subscriptions[i].cursor)
					if current < slowest {
						slowest = current
						subscribeOn = subscriptions[i].subscribeOn
					}
				}
				end := atomic.LoadUint64(&buf.end)
				if atomic.LoadUint64(&buf.begin) < slowest && slowest <= end {
					if slowest+buf.keep > end {
						slowest = end - buf.keep + 1
					}
					atomic.StoreUint64(&buf.begin, slowest)
					atomic.StoreUint64(&buf.end, slowest+buf.size)
				} else {
					if slowest == maxuint64 {
						atomic.AddUint64(&buf.begin, 1)
						atomic.AddUint64(&buf.end, 1)
					} else {
						full = true
					}
				}
			})
			if full {
				if !gosched {
					if subscribeOn != nil {
						subscribeOn.Gosched()
					} else {
						runtime.Gosched()
					}
				}
				if atomic.LoadUint64(&buf.state) != active {
					return
				}
			}
		}
		buf.items[buf.commit&buf.mod] = item{Value: value, At: time.Now()}
		atomic.AddUint64(&buf.commit, 1)
		buf.subscriptions.Broadcast()
	}

	close := func(err error) {
		if atomic.CompareAndSwapUint64(&buf.state, active, closing) {
			buf.err = err
			if atomic.CompareAndSwapUint64(&buf.state, closing, closed) {
				accessSubscriptions(func(subscriptions []subscription) {
					for i := range subscriptions {
						atomic.CompareAndSwapUint64(&subscriptions[i].state, active, closed)
					}
				})
			}
		}
		buf.subscriptions.Broadcast()
	}

	observer := func(next interface{}, err error, done bool) {
		if atomic.LoadUint64(&buf.state) == active {
			if !done {
				send(next)
			} else {
				close(err)
			}
		}
	}

	appendSubscription := func(subscribeOn Scheduler) (sub *subscription, err error) {
		accessSubscriptions(func([]subscription) {
			cursor := atomic.LoadUint64(&buf.begin)
			s := &buf.subscriptions
			if len(s.entries) < cap(s.entries) {
				s.entries = append(s.entries, subscription{cursor: cursor, subscribeOn: subscribeOn})
				sub = &s.entries[len(s.entries)-1]
				return
			}
			for i := range s.entries {
				sub = &s.entries[i]
				if atomic.CompareAndSwapUint64(&sub.cursor, maxuint64, cursor) {
					sub.subscribeOn = subscribeOn
					return
				}
			}
			sub = nil
			err = OutOfSubscriptions
			return
		})
		return
	}

	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		sub, err := appendSubscription(subscribeOn)
		if err != nil {
			runner := subscribeOn.Schedule(func() {
				if subscriber.Subscribed() {
					observe(nil, err, true)
				}
			})
			subscriber.OnUnsubscribe(runner.Cancel)
			return
		}
		commit := atomic.LoadUint64(&buf.commit)
		if atomic.LoadUint64(&buf.begin)+buf.keep < commit {
			atomic.StoreUint64(&sub.cursor, commit-buf.keep)
		}
		atomic.StoreUint64(&sub.state, atomic.LoadUint64(&buf.state))
		sub.activated = time.Now()

		receiver := subscribeOn.ScheduleFutureRecursive(0, func(self func(time.Duration)) {
			commit := atomic.LoadUint64(&buf.commit)

			if sub.cursor == commit {
				if atomic.CompareAndSwapUint64(&sub.state, canceled, canceled) {

					atomic.StoreUint64(&sub.cursor, maxuint64)
					return
				} else {

					now := time.Now()
					if now.Before(sub.activated.Add(1 * ms)) {

						self(50 * us)
						return
					} else if now.Before(sub.activated.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&sub.state, closed, closed) {

							observe(nil, buf.err, true)
							atomic.StoreUint64(&sub.cursor, maxuint64)
							return
						}

						self(500 * us)
						return
					} else {
						if subscribeOn.IsConcurrent() {

							buf.subscriptions.Lock()
							buf.subscriptions.Wait()
							buf.subscriptions.Unlock()
							sub.activated = time.Now()
							self(0)
							return
						} else {

							self(5 * ms)
							return
						}
					}
				}
			}

			if atomic.LoadUint64(&sub.state) == canceled {
				atomic.StoreUint64(&sub.cursor, maxuint64)
				return
			}
			for ; sub.cursor != commit; atomic.AddUint64(&sub.cursor, 1) {
				item := &buf.items[sub.cursor&buf.mod]
				if buf.age == 0 || item.At.IsZero() || time.Since(item.At) < buf.age {
					observe(item.Value, nil, false)
				}
				if atomic.LoadUint64(&sub.state) == canceled {
					atomic.StoreUint64(&sub.cursor, maxuint64)
					return
				}
			}

			sub.activated = time.Now()
			self(0)
		})
		subscriber.OnUnsubscribe(receiver.Cancel)

		subscriber.OnUnsubscribe(func() {
			atomic.CompareAndSwapUint64(&sub.state, active, canceled)
			buf.subscriptions.Broadcast()
		})
	}
	return observer, observable
}

//jig:name Subject

// Subject is a combination of an Observer and Observable.
// Subjects are special because they are the only reactive constructs that
// support multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// The Subject exposes all methods from the embedded Observer and
// Observable. Use the Observer Next, Error and Complete methods to feed
// data to it. Use the Observable methods to subscribe to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubject
// functions for more info.
type Subject struct {
	Observer
	Observable
}

// Next is called by an Observable to emit the next interface{} value to the
// Observer.
func (o Observer) Next(next interface{}) {
	o(next, nil, false)
}

// Error is called by an Observable to report an error to the Observer.
func (o Observer) Error(err error) {
	var zero interface{}
	o(zero, err, true)
}

// Complete is called by an Observable to signal that no more data is
// forthcoming to the Observer.
func (o Observer) Complete() {
	var zero interface{}
	o(zero, nil, true)
}

//jig:name Observer_AsObserver

// AsObserver converts an observer of interface{} items to an observer of
// interface{} items.
func (o Observer) AsObserver() Observer {
	observer := func(next interface{}, err error, done bool) {
		o(next, err, done)
	}
	return observer
}

//jig:name Observable_AsObservable

// AsObservable returns the source Observable unchanged.
// This is a special case needed for internal plumbing.
func (o Observable) AsObservable() Observable {
	return o
}

//jig:name NewAsyncSubject

// NewAsyncSubject creates a a subject that emits the last value (and only
// the last value) emitted by the Observable part, and only after that
// Observable part completes. (If the Observable part does not emit any
// values, the AsyncSubject also completes without emitting any values.)
//
// It will also emit this same final value to any subsequent observers.
// However, if the Observable part terminates with an error, the AsyncSubject
// will not emit any items, but will simply pass along the error notification
// from the Observable part.
func NewAsyncSubject() Subject {
	observer, observable := MakeObserverObservable(0, 1)
	var async struct {
		subject	Subject
		set	bool
		last	interface{}
	}
	async.subject.Observer = func(next interface{}, err error, done bool) {
		if !done {
			async.set = true
			async.last = next
		} else {
			if async.set && err == nil {
				observer(async.last, nil, false)
			}
			observer(next, err, true)
		}
	}
	async.subject.Observable = observable.AsObservable()
	return async.subject
}

//jig:name Observable_PublishLast

// PublishLast returns a Multicaster that shares a single subscription to
// the underlying Observable containing only the last value emitted before
// it completes. When the underlying Obervable terminates with an error,
// then subscribed observers will receive only that error (and no value).
// After all observers have unsubscribed due to an error, the Multicaster
// does an internal reset just before the next observer subscribes.
func (o Observable) PublishLast() Multicaster {
	return o.Multicast(NewAsyncSubject)
}

//jig:name ObservableInt_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
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

//jig:name ObservableInt_SubscribeOn

// SubscribeOn specifies the scheduler an ObservableInt should use when it is
// subscribed to.
func (o ObservableInt) SubscribeOn(scheduler Scheduler) ObservableInt {
	observable := func(observe IntObserver, _ Scheduler, subscriber Subscriber) {
		if scheduler.IsConcurrent() {
			subscriber.OnWait(nil)
		} else {
			subscriber.OnWait(scheduler.Wait)
		}
		o(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name Multicaster_RefCount

// RefCount makes a Multicaster behave like an ordinary Observable. On
// first Subscribe it will call Connect on its Multicaster and when its last
// subscriber is Unsubscribed it will cancel the source connection by calling
// Unsubscribe on the subscription returned by the call to Connect.
func (o Multicaster) RefCount() Observable {
	var source struct {
		sync.Mutex
		refcount	int32
		subscriber	Subscriber
	}
	observable := func(observe Observer, subscribeOn Scheduler, withSubscriber Subscriber) {
		withSubscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				source.subscriber.Unsubscribe()
			}
			source.Unlock()
		})
		o.Observable(observe, subscribeOn, withSubscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			source.subscriber = subscriber.New()
			source.Unlock()
			o.Connectable(subscribeOn, source.subscriber)
			source.Lock()
		}
		source.Unlock()
	}
	return observable
}

//jig:name Observable_Retry

// Retry if a source Observable sends an error notification, resubscribe to
// it in the hopes that it will complete without error. If count is zero or
// negative, the retry count will be effectively infinite. The scheduler
// passed when subscribing is used by Retry to schedule any retry attempt. The
// time between retries is 1 millisecond, so retry frequency is 1 kHz. Any
// SubscribeOn operators should be called after Retry to prevent lockups
// caused by mixing different schedulers in the same subscription for retrying
// and subscribing.
func (o Observable) Retry(count ...int) Observable {
	count = append(count, int(^uint(0)>>1))
	if count[0] <= 0 {
		count = []int{int(^uint(0) >> 1)}
	}
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var retry struct {
			count		int
			observer	Observer
			subscriber	Subscriber
			resubscribe	func()
		}
		retry.count = count[0]
		retry.observer = func(next interface{}, err error, done bool) {
			if err != nil && retry.count > 0 {
				retry.count--
				retry.subscriber.Done(err)
				subscribeOn.ScheduleFuture(1*time.Millisecond, retry.resubscribe)
			} else {
				observe(next, err, done)
			}
		}
		retry.resubscribe = func() {
			if subscriber.Subscribed() {
				retry.subscriber = subscriber.Add()
				if !subscribeOn.IsConcurrent() {
					retry.subscriber.OnWait(subscribeOn.Wait)
					subscriber.OnWait(func() { retry.subscriber.Wait() })
				}
				o(retry.observer, subscribeOn, retry.subscriber)
			}
		}
		retry.resubscribe()
	}
	return observable
}
