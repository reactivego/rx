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
		ms = time.Millisecond
		us = time.Microsecond
	)

	// Access to subscriptions
	const (
		idle uint32 = iota
		busy
	)

	// State of subscription and buffer
	const (
		active uint64 = iota
		canceled
		closing
		closed
	)

	const (
		// Cursor is parked so it does not influence advancing the commit index.
		parked uint64 = math.MaxUint64
	)

	type subscription struct {
		Cursor     uint64
		State      uint64    // active, canceled, closed
		LastActive time.Time // track activity to deterime backoff
	}

	type subscriptions struct {
		sync.Mutex
		*sync.Cond
		entries []subscription
		access  uint32 // idle, busy
	}

	type item struct {
		Value interface{}
		At    time.Time
	}

	type buffer struct {
		age time.Duration
		len uint64
		cap uint64

		mod    uint64
		items  []item
		begin  uint64
		end    uint64
		commit uint64
		state  uint64 // active, closed

		subscriptions subscriptions

		err error
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
			len:   len,
			cap:   cap,
			age:   age,
			items: make([]item, cap),
			mod:   cap - 1,
			end:   cap,
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
			})
			if atomic.LoadUint64(&buf.begin) < slowest && slowest <= atomic.LoadUint64(&buf.end) {
				atomic.StoreUint64(&buf.begin, slowest)
				atomic.StoreUint64(&buf.end, slowest+buf.mod+1)
			} else {
				slowest = parked
			}
			if slowest == parked {
				// no subscriptions present...
				if !spun {
					// if we did not spin to access subscriptions, then spin now.
					runtime.Gosched()
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

			// wait for commit to move away from cursor position, indicating fresh data
			if sub.Cursor == commit {
				if atomic.CompareAndSwapUint64(&sub.State, canceled, canceled) {
					// subscription canceled
					atomic.StoreUint64(&sub.Cursor, parked)
					return
				} else {
					// subscription still active (not canceled)
					now := time.Now()
					if now.Before(sub.LastActive.Add(1 * ms)) {
						// spinlock for 1ms (in increments of 50us) when no data from sender is arriving
						self(50 * us) // 20kHz
						return
					} else if now.Before(sub.LastActive.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&sub.State, closed, closed) {
							// buffer has been closed
							observe(nil, buf.err, true)
							atomic.StoreUint64(&sub.Cursor, parked)
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
							sub.LastActive = time.Now()
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

			// all caught up; record time and loop back to wait for fresh data
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
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer

// NewSubjectFoo creates a new Subject. After the subject is terminated, all
// subsequent subscriptions to the observable side will be terminated
// immediately with either an Error or Complete notification send to the
// subscribing client
//
// Note that this implementation is blocking. When there are subscribers, the
// observable goroutine is blocked until all subscribers have processed the
// next, error or complete notification.
func NewSubjectFoo() SubjectFoo {
	observer, observable := MakeObserverObservable(0, 0, 1, 16)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}

//jig:template MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:template NewReplaySubject<Foo>
//jig:needs MakeObserverObservable, Subject<Foo>, Observer As<Foo>Observer, MaxReplayCapacity

// NewReplaySubjectFoo creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
func NewReplaySubjectFoo(bufferCapacity int, windowDuration time.Duration) SubjectFoo {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	observer, observable := MakeObserverObservable(windowDuration, bufferCapacity)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}
