// Code generated by jig; DO NOT EDIT.

//go:generate jig

package AutoConnect

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name GoroutineScheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

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

//jig:name DeferInt

// DeferInt does not create the ObservableInt until the observer subscribes.
// It creates a fresh ObservableInt for each subscribing observer. Use it to
// create observables that maintain separate state per subscription.
func DeferInt(factory func() ObservableInt) ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name Error

// Error signals an error condition.
type Error func(error)

//jig:name Complete

// Complete signals that no more data is to be expected.
type Complete func()

//jig:name NextInt

// NextInt can be called to emit the next value to the IntObserver.
type NextInt func(int)

//jig:name CreateRecursiveInt

// CreateRecursiveInt provides a way of creating an ObservableInt from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateRecursiveInt will be called
// repeatedly to implement the observable. It is provided with a NextInt, Error
// and Complete function that can be called by the code that implements the
// Observable.
func CreateRecursiveInt(create func(NextInt, Error, Complete)) ObservableInt {
	var zeroInt int
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if !subscriber.Subscribed() {
				return
			}
			n := func(next int) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				done = true
				if subscriber.Subscribed() {
					observe(zeroInt, err, true)
				}
			}
			c := func() {
				done = true
				if subscriber.Subscribed() {
					observe(zeroInt, nil, true)
				}
			}
			create(n, e, c)
			if !done && subscriber.Subscribed() {
				self()
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Canceled

// Canceled returns true when the observer has unsubscribed.
type Canceled func() bool

//jig:name CreateInt

// CreateInt provides a way of creating an ObservableInt from
// scratch by calling observer methods programmatically.
//
// The create function provided to CreateInt will be called once
// to implement the observable. It is provided with a NextInt, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func CreateInt(create func(NextInt, Error, Complete, Canceled)) ObservableInt {
	var zeroInt int
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
					observe(zeroInt, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					observe(zeroInt, nil, true)
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

//jig:name FromInt

// FromInt creates an ObservableInt from multiple int values passed in.
func FromInt(slice ...int) ObservableInt {
	var zeroInt int
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
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
					observe(zeroInt, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ThrowInt

// ThrowInt creates an Observable that emits no items and terminates with an
// error.
func ThrowInt(err error) ObservableInt {
	var zeroInt int
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zeroInt, err, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name NeverInt

// NeverInt creates an ObservableInt that emits no items and does't terminate.
func NeverInt() ObservableInt {
	observable := func(observe IntObserver, scheduler Scheduler, subscriber Subscriber) {
	}
	return observable
}

//jig:name Timer

// Timer creates an ObservableInt that emits a sequence of integers (starting
// at zero) after an initialDelay has passed. Subsequent values are emitted
// using  a schedule of intervals passed in. If only the initialDelay is
// given, Timer will emit only once.
func Timer(initialDelay time.Duration, intervals ...time.Duration) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(initialDelay, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(i, nil, false)
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						self(intervals[i%len(intervals)])
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:name ObservableInt_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Subscribe(observe IntObserver, subscribers ...Subscriber) Subscription {
	subscribers = append(subscribers, subscriber.New())
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zeroInt int
			observe(zeroInt, err, true)
			subscribers[0].Done(err)
		}
	}
	subscribers[0].OnWait(scheduler.Wait)
	o(observer, scheduler, subscribers[0])
	return subscribers[0]
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

//jig:name IntMulticaster

// IntMulticaster is a multicasting connectable observable. One or more
// IntObservers can subscribe to it simultaneously. It will subscribe to the
// source ObservableInt when Connect is called. After that, every emission
// from the source is multcast to all subscribed IntObservers.
type IntMulticaster struct {
	ObservableInt
	Connectable
}

//jig:name ObservableInt_Multicast

// Multicast converts an ordinary observable into a multicasting connectable
// observable or multicaster for short. A multicaster will only start emitting
// values after its Connect method has been called. The factory method passed
// in should return a new SubjectInt that implements the actual multicasting
// behavior.
func (o ObservableInt) Multicast(factory func() SubjectInt) IntMulticaster {
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
	observer := func(next int, err error, done bool) {
		if atomic.CompareAndSwapInt32(&subject.state, active, notifying) {
			if s, ok := subject.Load().(SubjectInt); ok {
				s.IntObserver(next, err, done)
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
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		if atomic.AddInt32(&subject.count, 1) == 1 {
			if atomic.CompareAndSwapInt32(&subject.state, erred, active) {
				subject.Store(factory())
			}
		}
		if s, ok := subject.Load().(SubjectInt); ok {
			s.ObservableInt(observe, subscribeOn, subscriber)
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
	return IntMulticaster{ObservableInt: observable, Connectable: connectable}
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

//jig:name MakeObserverObservable

const ErrOutOfSubscriptions = RxError("out of subscriptions")

// MakeObserverObservable actually does make an observer observable. It
// creates a buffering multicaster and returns both the Observer and the
// Observable side of it. These are then used as the core of any Subject
// implementation. The Observer side is used to pass items into the buffering
// multicaster. This then multicasts the items to every Observer that
// subscribes to the returned Observable.
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

	// Access to subscriptions
	const (
		idle	uint32	= iota
		busy
	)

	// State of subscription and buffer
	const (
		active	uint64	= iota
		canceled
		closing
		closed
	)

	const (
		// Cursor is parked so it does not influence advancing the commit index.
		parked uint64 = math.MaxUint64
	)

	type subscription struct {
		Cursor		uint64
		State		uint64		// active, canceled, closed
		LastActive	time.Time	// track activity to deterime backoff
	}

	type subscriptions struct {
		sync.Mutex
		*sync.Cond
		entries	[]subscription
		access	uint32	// idle, busy
	}

	type item struct {
		Value	interface{}
		At	time.Time
	}

	type buffer struct {
		age	time.Duration
		len	uint64
		cap	uint64

		mod	uint64
		items	[]item
		begin	uint64
		end	uint64
		commit	uint64
		state	uint64	// active, closed

		subscriptions	subscriptions

		err	error
	}

	make := func(age time.Duration, length int, capacity ...int) *buffer {
		cap, scap := uint64(length), uint64(32)
		switch {
		case len(capacity) >= 2:
			cap, scap = uint64(capacity[0]), uint64(capacity[1])
		case len(capacity) == 1:
			cap = uint64(capacity[0])
		}
		len := uint64(length)
		if cap < len {
			cap = len
		}
		cap = uint64(1) << uint(math.Ceil(math.Log2(float64(cap))))
		buf := &buffer{
			len:	len,
			cap:	cap,
			age:	age,
			items:	make([]item, cap),
			mod:	cap - 1,
			end:	cap,
			subscriptions: subscriptions{
				entries: make([]subscription, 0, scap),
			},
		}
		buf.subscriptions.Cond = sync.NewCond(&buf.subscriptions.Mutex)
		return buf
	}
	buf := make(age, length, capacity...)

	accessSubscriptions := func(access func([]subscription)) bool {
		spun := false
		for !atomic.CompareAndSwapUint32(&buf.subscriptions.access, idle, busy) {
			runtime.Gosched()
			spun = true
		}
		access(buf.subscriptions.entries)
		atomic.StoreUint32(&buf.subscriptions.access, idle)
		return spun
	}

	send := func(value interface{}) {
		for buf.commit == buf.end {
			slowest := parked
			spun := accessSubscriptions(func(subscriptions []subscription) {
				for i := range subscriptions {
					cursor := atomic.LoadUint64(&subscriptions[i].Cursor)
					if cursor < slowest {
						slowest = cursor
					}
				}
				if atomic.LoadUint64(&buf.begin) < slowest && slowest <= atomic.LoadUint64(&buf.end) {
					atomic.StoreUint64(&buf.begin, slowest)
					atomic.StoreUint64(&buf.end, slowest+buf.mod+1)
				} else {
					slowest = parked
				}
			})
			if slowest == parked {

				if !spun {

					runtime.Gosched()
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
						atomic.CompareAndSwapUint64(&subscriptions[i].State, active, closed)
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

	appendSubscription := func(cursor uint64) (sub *subscription, err error) {
		accessSubscriptions(func([]subscription) {
			s := &buf.subscriptions
			if len(s.entries) < cap(s.entries) {
				s.entries = append(s.entries, subscription{Cursor: cursor})
				sub = &s.entries[len(s.entries)-1]
				return
			}
			for i := range s.entries {
				sub = &s.entries[i]
				if atomic.CompareAndSwapUint64(&sub.Cursor, parked, cursor) {
					return
				}
			}
			err = ErrOutOfSubscriptions
			return
		})
		return
	}

	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		begin := atomic.LoadUint64(&buf.begin)
		sub, err := appendSubscription(begin)
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
		if begin+buf.len < commit {
			atomic.StoreUint64(&sub.Cursor, commit-buf.len)
		}
		atomic.StoreUint64(&sub.State, atomic.LoadUint64(&buf.state))
		sub.LastActive = time.Now()

		receiver := subscribeOn.ScheduleFutureRecursive(0, func(self func(time.Duration)) {
			commit := atomic.LoadUint64(&buf.commit)

			if sub.Cursor == commit {
				if atomic.CompareAndSwapUint64(&sub.State, canceled, canceled) {

					atomic.StoreUint64(&sub.Cursor, parked)
					return
				} else {

					now := time.Now()
					if now.Before(sub.LastActive.Add(1 * ms)) {

						self(50 * us)
						return
					} else if now.Before(sub.LastActive.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&sub.State, closed, closed) {

							observe(nil, buf.err, true)
							atomic.StoreUint64(&sub.Cursor, parked)
							return
						}

						self(500 * us)
						return
					} else {
						if subscribeOn.IsConcurrent() {

							buf.subscriptions.Lock()
							buf.subscriptions.Wait()
							buf.subscriptions.Unlock()
							sub.LastActive = time.Now()
							self(0)
							return
						} else {

							self(5 * ms)
							return
						}
					}
				}
			}

			if atomic.LoadUint64(&sub.State) == canceled {
				atomic.StoreUint64(&sub.Cursor, parked)
				return
			}
			for ; sub.Cursor != commit; atomic.AddUint64(&sub.Cursor, 1) {
				item := &buf.items[sub.Cursor&buf.mod]
				if buf.age == 0 || item.At.IsZero() || time.Since(item.At) < buf.age {
					observe(item.Value, nil, false)
				}
				if atomic.LoadUint64(&sub.State) == canceled {
					atomic.StoreUint64(&sub.Cursor, parked)
					return
				}
			}

			sub.LastActive = time.Now()
			self(0)
		})
		subscriber.OnUnsubscribe(receiver.Cancel)

		subscriber.OnUnsubscribe(func() {
			atomic.CompareAndSwapUint64(&sub.State, active, canceled)
			buf.subscriptions.Broadcast()
		})
	}
	return observer, observable
}

//jig:name SubjectInt

// SubjectInt is a combination of an IntObserver and ObservableInt.
// Subjects are special because they are the only reactive constructs that
// support multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// The SubjectInt exposes all methods from the embedded IntObserver and
// ObservableInt. Use the IntObserver Next, Error and Complete methods to feed
// data to it. Use the ObservableInt methods to subscribe to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubjectInt
// functions for more info.
type SubjectInt struct {
	IntObserver
	ObservableInt
}

// Next is called by an ObservableInt to emit the next int value to the
// Observer.
func (o IntObserver) Next(next int) {
	o(next, nil, false)
}

// Error is called by an ObservableInt to report an error to the Observer.
func (o IntObserver) Error(err error) {
	var zero int
	o(zero, err, true)
}

// Complete is called by an ObservableInt to signal that no more data is
// forthcoming to the Observer.
func (o IntObserver) Complete() {
	var zero int
	o(zero, nil, true)
}

//jig:name Observer_AsIntObserver

// AsIntObserver converts an observer of interface{} items to an observer of
// int items.
func (o Observer) AsIntObserver() IntObserver {
	observer := func(next int, err error, done bool) {
		o(next, err, done)
	}
	return observer
}

//jig:name MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:name NewReplaySubjectInt

// NewReplaySubjectInt creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
func NewReplaySubjectInt(bufferCapacity int, windowDuration time.Duration) SubjectInt {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	observer, observable := MakeObserverObservable(windowDuration, bufferCapacity)
	return SubjectInt{observer.AsIntObserver(), observable.AsObservableInt()}
}

//jig:name ObservableInt_PublishReplay

// Replay uses the Multicast operator to control the subscription of a
// ReplaySubject to a source observable and turns the subject into a
// connectable observable. A ReplaySubject emits to any observer all of the
// items that were emitted by the source observable, regardless of when the
// observer subscribes.
//
// If the source completed and as a result the internal ReplaySubject
// terminated, then calling Connect again will replace the old ReplaySubject
// with a newly created one.
func (o ObservableInt) PublishReplay(bufferCapacity int, windowDuration time.Duration) IntMulticaster {
	factory := func() SubjectInt {
		return NewReplaySubjectInt(bufferCapacity, windowDuration)
	}
	return o.Multicast(factory)
}

//jig:name NewSubjectInt

// NewSubjectInt creates a new Subject. After the subject is terminated, all
// subsequent subscriptions to the observable side will be terminated
// immediately with either an Error or Complete notification send to the
// subscribing client
//
// Note that this implementation is blocking. When there are subscribers, the
// observable goroutine is blocked until all subscribers have processed the
// next, error or complete notification.
func NewSubjectInt() SubjectInt {
	observer, observable := MakeObserverObservable(0, 0, 1, 16)
	return SubjectInt{observer.AsIntObserver(), observable.AsObservableInt()}
}

//jig:name ObservableInt_Publish

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
func (o ObservableInt) Publish() IntMulticaster {
	return o.Multicast(NewSubjectInt)
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
const TimeoutOccured = RxError("timeout")

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

//jig:name ObservableInt_Timeout

// Timeout mirrors the source ObservableInt, but issues an error notification if
// a particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o ObservableInt) Timeout(timeout time.Duration) ObservableInt {
	return o.AsObservable().Timeout(timeout).AsObservableInt()
}

//jig:name ErrTypecastToInt

// ErrTypecastToInt is delivered to an observer if the generic value cannot be
// typecast to int.
const ErrTypecastToInt = RxError("typecast to int failed")

//jig:name Observable_AsObservableInt

// AsObservableInt turns an Observable of interface{} into an ObservableInt.
// If during observing a typecast fails, the error ErrTypecastToInt will be
// emitted.
func (o Observable) AsObservableInt() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextInt, ok := next.(int); ok {
					observe(nextInt, err, done)
				} else {
					var zeroInt int
					observe(zeroInt, ErrTypecastToInt, true)
				}
			} else {
				var zeroInt int
				observe(zeroInt, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableInt_AsObservable

// AsObservable turns a typed ObservableInt into an Observable of interface{}.
func (o ObservableInt) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next int, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ErrAutoConnect

const InvalidCount = RxError("invalid count")

//jig:name IntMulticaster_AutoConnect

// AutoConnect makes a IntMulticaster behave like an ordinary ObservableInt
// that automatically connects the multicaster to its source when the
// specified number of observers have subscribed to it. If the count is less
// than 1 it will return a ThrowInt(InvalidCount). After
// connecting, when the number of subscribed observers eventually drops to 0,
// AutoConnect will cancel the source connection if it hasn't terminated yet.
// When subsequently the next observer subscribes, AutoConnect will connect to
// the source only when it was previously canceled or because the source
// terminated with an error. So it will not reconnect when the source
// completed succesfully. This specific behavior allows for implementing a
// caching observable that can be retried until it succeeds. Another thing to
// notice is that AutoConnect will disconnect an active connection when the
// number of observers drops to zero. The reason for this is that not doing so
// would leak a task and leave it hanging in the scheduler.
func (o IntMulticaster) AutoConnect(count int) ObservableInt {
	if count < 1 {
		return ThrowInt(InvalidCount)
	}
	var source struct {
		sync.Mutex
		refcount	int32
		subscriber	Subscriber
	}
	observable := func(observe IntObserver, subscribeOn Scheduler, withSubscriber Subscriber) {
		withSubscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				if source.subscriber != nil {
					source.subscriber.Unsubscribe()
				}
			}
			source.Unlock()
		})
		o.ObservableInt(observe, subscribeOn, withSubscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == int32(count) {
			if source.subscriber == nil || source.subscriber.Error() != nil {
				source.subscriber = subscriber.New()
				source.Unlock()
				o.Connectable(subscribeOn, source.subscriber)
				source.Lock()
			}
		}
		source.Unlock()
	}
	return observable
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

//jig:name ObservableInt_Wait

// Wait subscribes to the Observable and waits for completion or error.
// Returns either the error or nil when the Observable completed normally.
// Wait uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) Wait() error {
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
		if done {
			subscriber.Done(err)
		}
	}
	subscriber.OnWait(scheduler.Wait)
	o(observer, scheduler, subscriber)
	return subscriber.Wait()
}

//jig:name ObservableInt_ToSingle

// ToSingle blocks until the ObservableInt emits exactly one value or an error.
// The value and any error are returned.
// ToSingle uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o ObservableInt) ToSingle() (entry int, err error) {
	o = o.Single()
	subscriber := subscriber.New()
	scheduler := scheduler.MakeTrampoline()
	observer := func(next int, err error, done bool) {
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

//jig:name ObservableInt_Single

// Single enforces that the observableInt sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o ObservableInt) Single() ObservableInt {
	return o.AsObservable().Single().AsObservableInt()
}

//jig:name ErrSingle

const DidNotEmitValue = RxError("expected one value, got none")

const EmittedMultipleValues = RxError("expected one value, got multiple")

//jig:name Observable_Single

// Single enforces that the observable sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing, this is reported as an error to the observer.
func (o Observable) Single() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			count	int
			latest	interface{}
		)
		observer := func(next interface{}, err error, done bool) {
			if count < 2 {
				if done {
					if err != nil {
						observe(nil, err, true)
					} else {
						if count == 1 {
							observe(latest, nil, false)
							observe(nil, nil, true)
						} else {
							observe(nil, DidNotEmitValue, true)
						}
					}
				} else {
					count++
					if count == 1 {
						latest = next
					} else {
						observe(nil, EmittedMultipleValues, true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
