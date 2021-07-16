package rx

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

//jig:template MakeObserverObservable
//jig:needs Observer, Observable

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
		ms = time.Millisecond
		us = time.Microsecond
	)

	type subscription struct {
		cursor      uint64
		state       uint64    // active, canceled, closed
		activated   time.Time // track activity to deterime backoff
		subscribeOn Scheduler
	}

	// cursor
	const (
		maxuint64 uint64 = math.MaxUint64 // park unused cursor at maxuint64
	)

	// state
	const (
		active uint64 = iota
		canceled
		closing
		closed
	)

	type subscriptions struct {
		sync.Mutex
		*sync.Cond
		entries []subscription
		access  uint32 // unlocked, locked
	}

	// access
	const (
		unlocked uint32 = iota
		locked
	)

	type item struct {
		Value interface{}
		At    time.Time
	}

	type buffer struct {
		age  time.Duration
		keep uint64
		mod  uint64
		size uint64

		items  []item
		begin  uint64
		end    uint64
		commit uint64
		state  uint64 // active, closed

		subscriptions subscriptions

		err error
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
		size := uint64(1) << uint(math.Ceil(math.Log2(float64(cap)))) // MUST(keep < size)

		if scap < 1 {
			scap = 1
		}
		buf := &buffer{
			age:  age,
			keep: keep,
			mod:  size - 1,
			size: size,

			items: make([]item, size),
			end:   size,
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
					if slowest == maxuint64 { // no subscriptions
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
					return // buffer no longer active
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

			// wait for commit to move away from cursor position, indicating fresh data
			if sub.cursor == commit {
				if atomic.CompareAndSwapUint64(&sub.state, canceled, canceled) {
					// subscription canceled
					atomic.StoreUint64(&sub.cursor, maxuint64)
					return
				} else {
					// subscription still active (not canceled)
					now := time.Now()
					if now.Before(sub.activated.Add(1 * ms)) {
						// spinlock for 1ms (in increments of 50us) when no data from sender is arriving
						self(50 * us) // 20kHz
						return
					} else if now.Before(sub.activated.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&sub.state, closed, closed) {
							// buffer has been closed
							observe(nil, buf.err, true)
							atomic.StoreUint64(&sub.cursor, maxuint64)
							return
						}
						// spinlock between 1ms and 250ms (in increments of 500us) of no data from sender
						self(500 * us) // 2kHz
						return
					} else {
						if subscribeOn.IsConcurrent() {
							// Block goroutine on condition until notified.
							buf.subscriptions.Lock()
							buf.subscriptions.Wait()
							buf.subscriptions.Unlock()
							sub.activated = time.Now()
							self(0)
							return
						} else {
							// Spinlock rest of the time (in increments of 5ms). This wakes-up
							// slower than the condition based solution above, but can be used with
							// a trampoline scheduler.
							self(5 * ms) // 200Hz
							return
						}
					}
				}
			}

			// emit data and advance cursor to catch up to commit
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

			// all caught up; record time and loop back to wait for fresh data
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

//jig:template Subject<Foo>
//jig:embeds <Foo>Observer, Observable<Foo>

// SubjectFoo is a combination of an FooObserver and ObservableFoo.
// Subjects are special because they are the only reactive constructs that
// support multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// The SubjectFoo exposes all methods from the embedded FooObserver and
// ObservableFoo. Use the FooObserver Next, Error and Complete methods to feed
// data to it. Use the ObservableFoo methods to subscribe to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubjectFoo
// functions for more info.
type SubjectFoo struct {
	FooObserver
	ObservableFoo
}

// Next is called by an ObservableFoo to emit the next foo value to the
// Observer.
func (o FooObserver) Next(next foo) {
	o(next, nil, false)
}

// Error is called by an ObservableFoo to report an error to the Observer.
func (o FooObserver) Error(err error) {
	var zero foo
	o(zero, err, true)
}

// Complete is called by an ObservableFoo to signal that no more data is
// forthcoming to the Observer.
func (o FooObserver) Complete() {
	var zero foo
	o(zero, nil, true)
}

//jig:template NewSubject<Foo>
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer, Observable AsObservable<Foo>

// NewSubjectFoo creates a new Subject. After the subject is terminated, all
// subsequent subscriptions to the observable side will be terminated
// immediately with either an Error or Complete notification send to the
// subscribing client
//
// Note that this implementation is blocking. When there are subscribers, the
// observable goroutine is blocked until all subscribers have processed the
// next, error or complete notification.
func NewSubjectFoo() SubjectFoo {
	observer, observable := MakeObserverObservable(0, 0)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}

//jig:template DefaultReplayCapacity

// DefaultReplayCapacity is the default capacity of a replay buffer when
// a bufferCapacity of 0 is passed to the NewReplaySubject function.
const DefaultReplayCapacity = 16383

//jig:template NewReplaySubject<Foo>
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer, Observable AsObservable<Foo>, DefaultReplayCapacity

// NewReplaySubjectFoo creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then DefaultReplayCapacity is
// used (currently 16380). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
func NewReplaySubjectFoo(bufferCapacity int, windowDuration time.Duration) SubjectFoo {
	if bufferCapacity == 0 {
		bufferCapacity = DefaultReplayCapacity
	}
	observer, observable := MakeObserverObservable(windowDuration, bufferCapacity)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}

//jig:template NewBehaviorSubject<Foo>
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer, Observable AsObservable<Foo>

// NewBehaviorSubjectFoo returns a new BehaviorSubjectFoo. When an observer
// subscribes to a BehaviorSubject, it begins by emitting the item most
// recently emitted by the Observable part of the subject (or a seed/default
// value if none has yet been emitted) and then continues to emit any other
// items emitted later by the Observable part.
func NewBehaviorSubjectFoo(a foo) SubjectFoo {
	observer, observable := MakeObserverObservable(0, 1)
	var behavior struct {
		subject   SubjectFoo
		completed uint32
	}
	observer(a, nil, false)
	behavior.subject.FooObserver = func(next foo, err error, done bool) {
		if done && err == nil {
			atomic.StoreUint32(&behavior.completed, 1)
		}
		observer(next, err, done)
	}
	behavior.subject.ObservableFoo = func(observe FooObserver, subscribeOn Scheduler, subscribe Subscriber) {
		if atomic.LoadUint32(&behavior.completed) != 1 {
			o := observable.AsObservableFoo()
			o(observe, subscribeOn, subscribe)
		} else {
			o := EmptyFoo()
			o(observe, subscribeOn, subscribe)
		}
	}
	return behavior.subject
}

//jig:template NewAsyncSubject<Foo>
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer, Observable AsObservable<Foo>

// NewAsyncSubjectFoo creates a a subject that emits the last value (and only
// the last value) emitted by the Observable part, and only after that
// Observable part completes. (If the Observable part does not emit any
// values, the AsyncSubject also completes without emitting any values.)
//
// It will also emit this same final value to any subsequent observers.
// However, if the Observable part terminates with an error, the AsyncSubject
// will not emit any items, but will simply pass along the error notification
// from the Observable part.
func NewAsyncSubjectFoo() SubjectFoo {
	observer, observable := MakeObserverObservable(0, 1)
	var async struct {
		subject SubjectFoo
		set     bool
		last    foo
	}
	async.subject.FooObserver = func(next foo, err error, done bool) {
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
	async.subject.ObservableFoo = observable.AsObservableFoo()
	return async.subject
}
