package rx

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const OutOfSubjectSubscriptions = Error("out of subject subscriptions")

// Subject returns both an Observer and and Observable. The returned Observer is
// used to send items into the Subject. The returned Observable is used to subscribe
// to the Subject. The Subject multicasts items send through the Observer to every
// Subscriber of the Observable.
//
//	age     max age to keep items in order to replay them to a new Subscriber (0 = no max age).
//	[size]  size of the item buffer, number of items kept to replay to a new Subscriber.
//	[cap]   capacity of the item buffer, number of items that can be observed before blocking.
//	[scap]  capacity of the subscription list, max number of simultaneous subscribers.
func Subject[T any](age time.Duration, capacity ...int) (Observer[T], Observable[T]) {
	const (
		ms = time.Millisecond
		us = time.Microsecond
	)

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

	// access
	const (
		unlocked uint32 = iota
		locked
	)

	make := func(age time.Duration, capacity ...int) *_buffer[T] {
		size, cap, scap := 0, 0, 32
		switch {
		case len(capacity) == 3:
			size, cap, scap = capacity[0], capacity[1], capacity[2]
		case len(capacity) == 2:
			size, cap = capacity[0], capacity[1]
		case len(capacity) == 1:
			size = capacity[0]
		}
		if size < 0 {
			size = 0
		}
		if cap < size {
			cap = size
		}
		if cap == 0 {
			cap = 1
		}
		length := uint64(1) << uint(math.Ceil(math.Log2(float64(cap)))) // MUST(keep < length)

		if scap < 1 {
			scap = 1
		}
		buf := &_buffer[T]{
			age:  age,
			keep: uint64(size),
			mod:  length - 1,
			size: length,

			items: make([]_item[T], length),
			end:   length,
			subscriptions: _subscriptions{
				entries: make([]_subscription, 0, scap),
			},
		}
		buf.subscriptions.Cond = sync.NewCond(&buf.subscriptions.Mutex)
		return buf
	}
	buf := make(age, capacity...)

	accessSubscriptions := func(access func([]_subscription)) bool {
		gosched := false
		for !atomic.CompareAndSwapUint32(&buf.subscriptions.access, unlocked, locked) {
			runtime.Gosched()
			gosched = true
		}
		access(buf.subscriptions.entries)
		atomic.StoreUint32(&buf.subscriptions.access, unlocked)
		return gosched
	}

	send := func(value T) {
		for buf.commit == buf.end {
			full := false
			scheduler := Scheduler(nil)
			gosched := accessSubscriptions(func(subscriptions []_subscription) {
				slowest := maxuint64
				for i := range subscriptions {
					current := atomic.LoadUint64(&subscriptions[i].cursor)
					if current < slowest {
						slowest = current
						scheduler = subscriptions[i].scheduler
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
					if scheduler != nil {
						scheduler.Gosched()
					} else {
						runtime.Gosched()
					}
				}
				if atomic.LoadUint64(&buf.state) != active {
					return // buffer no longer active
				}
			}
		}
		buf.items[buf.commit&buf.mod] = _item[T]{Value: value, At: time.Now()}
		atomic.AddUint64(&buf.commit, 1)
		buf.subscriptions.Broadcast()
	}

	close := func(err error) {
		if atomic.CompareAndSwapUint64(&buf.state, active, closing) {
			buf.err = err
			if atomic.CompareAndSwapUint64(&buf.state, closing, closed) {
				accessSubscriptions(func(subscriptions []_subscription) {
					for i := range subscriptions {
						atomic.CompareAndSwapUint64(&subscriptions[i].state, active, closed)
					}
				})
			}
		}
		buf.subscriptions.Broadcast()
	}

	observer := func(next T, err error, done bool) {
		if atomic.LoadUint64(&buf.state) == active {
			if !done {
				send(next)
			} else {
				close(err)
			}
		}
	}

	appendSubscription := func(scheduler Scheduler) (sub *_subscription, err error) {
		accessSubscriptions(func([]_subscription) {
			cursor := atomic.LoadUint64(&buf.begin)
			s := &buf.subscriptions
			if len(s.entries) < cap(s.entries) {
				s.entries = append(s.entries, _subscription{cursor: cursor, scheduler: scheduler})
				sub = &s.entries[len(s.entries)-1]
				return
			}
			for i := range s.entries {
				sub = &s.entries[i]
				if atomic.CompareAndSwapUint64(&sub.cursor, maxuint64, cursor) {
					sub.scheduler = scheduler
					return
				}
			}
			sub = nil
			err = OutOfSubjectSubscriptions
		})
		return
	}

	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		sub, err := appendSubscription(scheduler)
		if err != nil {
			runner := scheduler.Schedule(func() {
				if subscriber.Subscribed() {
					var zero T
					observe(zero, err, true)
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

		receiver := scheduler.ScheduleFutureRecursive(0, func(self func(time.Duration)) {
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
							var zero T
							observe(zero, buf.err, true)
							atomic.StoreUint64(&sub.cursor, maxuint64)
							return
						}
						// spinlock between 1ms and 250ms (in increments of 500us) of no data from sender
						self(500 * us) // 2kHz
						return
					} else {
						if scheduler.IsConcurrent() {
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

type _buffer[T any] struct {
	age  time.Duration
	keep uint64
	mod  uint64
	size uint64

	items  []_item[T]
	begin  uint64
	end    uint64
	commit uint64
	state  uint64 // active, closed

	subscriptions _subscriptions

	err error
}

type _item[T any] struct {
	Value T
	At    time.Time
}

type _subscriptions struct {
	sync.Mutex
	*sync.Cond
	entries []_subscription
	access  uint32 // unlocked, locked
}

type _subscription struct {
	cursor    uint64
	state     uint64    // active, canceled, closed
	activated time.Time // track activity to deterime backoff
	scheduler Scheduler
}
