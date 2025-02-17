// Code generated by jig; DO NOT EDIT.

//go:generate jig

package PublishBehavior

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

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

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

//jig:name GoroutineScheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

//jig:name MakeTrampolineScheduler

func MakeTrampolineScheduler() Scheduler {
	return scheduler.MakeTrampoline()
}

//jig:name Concat

// Concat emits the emissions from two or more Observables without interleaving them.
func Concat(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return observables[0].ConcatWith(observables[1:]...)
}

//jig:name Timer

// Timer creates an Observable that emits a sequence of integers
// (starting at zero) after an initialDelay has passed. Subsequent values are
// emitted using  a schedule of intervals passed in. If only the initialDelay
// is given, Timer will emit only once.
func Timer(initialDelay time.Duration, intervals ...time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(initialDelay, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				if i == 0 || (i > 0 && len(intervals) > 0) {
					observe(interface{}(i), nil, false)
				}
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						self(intervals[i%len(intervals)])
					} else {
						if i == 0 {
							self(0)
						} else {
							var zero interface{}
							observe(zero, nil, true)
						}
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Just

// Just creates an Observable that emits a particular item.
func Just(element interface{}) Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(element, nil, false)
			}
			if subscriber.Subscribed() {
				var zero interface{}
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

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

//jig:name Empty

// Empty creates an Observable that emits no items but terminates normally.
func Empty() Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				var zero interface{}
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Observable_ConcatWith

// ConcatWith emits the emissions from two or more Observables without interleaving them.
func (o Observable) ConcatWith(other ...Observable) Observable {
	if len(other) == 0 {
		return o
	}
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			observables	= append([]Observable{}, other...)
			observer	Observer
		)
		observer = func(next interface{}, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					var zero interface{}
					observe(zero, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
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

//jig:name NewBehaviorSubject

// NewBehaviorSubject returns a new BehaviorSubject. When an observer
// subscribes to a BehaviorSubject, it begins by emitting the item most
// recently emitted by the Observable part of the subject (or a seed/default
// value if none has yet been emitted) and then continues to emit any other
// items emitted later by the Observable part.
func NewBehaviorSubject(a interface{}) Subject {
	observer, observable := MakeObserverObservable(0, 1)
	var behavior struct {
		subject		Subject
		completed	uint32
	}
	observer(a, nil, false)
	behavior.subject.Observer = func(next interface{}, err error, done bool) {
		if done && err == nil {
			atomic.StoreUint32(&behavior.completed, 1)
		}
		observer(next, err, done)
	}
	behavior.subject.Observable = func(observe Observer, subscribeOn Scheduler, subscribe Subscriber) {
		if atomic.LoadUint32(&behavior.completed) != 1 {
			o := observable.AsObservable()
			o(observe, subscribeOn, subscribe)
		} else {
			o := Empty()
			o(observe, subscribeOn, subscribe)
		}
	}
	return behavior.subject
}

//jig:name Observable_PublishBehavior

// PublishBehavior returns a Multicaster that shares a single subscription
// to the underlying Observable returning an initial value or the last
// value emitted by the underlying Observable. When the underlying
// Obervable terminates with an error, then subscribed observers will
// receive that error. After all observers have unsubscribed due to an error,
// the Multicaster does an internal reset just before the next observer
// subscribes.
func (o Observable) PublishBehavior(a interface{}) Multicaster {
	factory := func() Subject {
		return NewBehaviorSubject(a)
	}
	return o.Multicast(factory)
}

//jig:name Observable_MergeMapTo

// MergeMapTo maps every entry emitted by the Observable into a single
// Observable. The stream of Observable items is then merged into a
// single stream of  items using the MergeAll operator.
func (o Observable) MergeMapTo(inner Observable) Observable {
	project := func(interface{}) Observable { return inner }
	return o.MapObservable(project).MergeAll()
}

//jig:name Observable_MapObservable

// MapObservable transforms the items emitted by an Observable by applying a
// function to each item.
func (o Observable) MapObservable(project func(interface{}) Observable) ObservableObservable {
	observable := func(observe ObservableObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			var mapped Observable
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_Println

// Println subscribes to the Observable and prints every item to os.Stdout
// while it waits for completion or error. Returns either the error or nil
// when the Observable completed normally.
// Println uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o Observable) Println(a ...interface{}) error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next interface{}, err error, done bool) {
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

//jig:name Observable_Serialize

// Serialize forces an Observable to make serialized calls and to be
// well-behaved.
func (o Observable) Serialize() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var observer struct {
			sync.Mutex
			done	bool
		}
		serializer := func(next interface{}, err error, done bool) {
			observer.Lock()
			defer observer.Unlock()
			if !observer.done {
				observer.done = done
				observe(next, err, done)
			}
		}
		o(serializer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_Timeout

// TimeoutOccured is delivered to an observer if the stream times out.
const TimeoutOccured = RxError("timeout occured")

// Timeout mirrors the source Observable, but issues an error notification if a
// particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o Observable) Timeout(due time.Duration) Observable {
	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var timeout struct {
			sync.Mutex
			at		time.Time
			occurred	bool
		}
		timeout.at = subscribeOn.Now().Add(due)
		timer := subscribeOn.ScheduleFutureRecursive(due, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					due := timeout.at.Sub(subscribeOn.Now())
					if due > 0 {
						self(due)
					} else {
						timeout.occurred = true
						timeout.Unlock()
						observe(nil, TimeoutOccured, true)
						timeout.Lock()
					}
				}
				timeout.Unlock()
			}
		})
		subscriber.OnUnsubscribe(timer.Cancel)
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					now := subscribeOn.Now()
					if now.Before(timeout.at) {
						timeout.at = now.Add(due)
						timeout.occurred = done
						observe(next, err, done)
					}
				}
				timeout.Unlock()
			}
		}
		o(observer, subscribeOn, subscriber)
	})
	return observable
}

//jig:name ObservableObserver

// ObservableObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableObserver func(next Observable, err error, done bool)

//jig:name ObservableObservable

// ObservableObservable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservable func(ObservableObserver, Scheduler, Subscriber)

//jig:name ObservableObservable_MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservable) MergeAll() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done	bool
			len	int32
		}
		observer := func(next interface{}, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					var zero interface{}
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zero interface{}
						observe(zero, nil, true)
					}
				}
			}
		}
		merger := func(next Observable, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next.AutoUnsubscribe()(observer, subscribeOn, subscriber)
			} else {
				var zero interface{}
				observer(zero, err, true)
			}
		}
		observers.len += 1
		o.AutoUnsubscribe()(merger, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Observable_AutoUnsubscribe

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o Observable) AutoUnsubscribe() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next interface{}, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableObservable_AutoUnsubscribeObservable

// AutoUnsubscribe will automatically unsubscribe from the source when it signals it is done.
// This Operator subscribes to the source Observable using a separate subscriber. When the source
// observable subsequently signals it is done, the separate subscriber will be Unsubscribed.
func (o ObservableObservable) AutoUnsubscribe() ObservableObservable {
	observable := func(observe ObservableObserver, subscribeOn Scheduler, subscriber Subscriber) {
		subscriber = subscriber.Add()
		observer := func(next Observable, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
