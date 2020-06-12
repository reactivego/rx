package rx

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

//jig:template NewBuffer
//jig:needs Observer, Observable

const ErrOutOfEndpoints = RxError("out of endpoints")

// NewBuffer creates a buffer to be used as the core of any Subject
// implementation. It returns both an Observer as well as an Observable. Items
// are placed in the buffer through the returned Observer. The buffer then
// multicasts the item to every subscriber of the returned Observable.
// 
//	age     age below which items are kept to replay to a new subscriber.
//	length  length of the item buffer, number of items kept to replay to a new subscriber.
//	[cap]   Capacity of the item buffer, number of items that can be observed before blocking.
//	[ecap]  Capacity of the endpoints slice.
func NewBuffer(age time.Duration, length int, capacity ...int) (Observer, Observable) {
	const (
		ms = time.Millisecond
		us = time.Microsecond
	)

	// Access to endpoints
	const (
		idle uint32 = iota
		busy
	)

	// State of endpoint and Chan
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

	type endpoint struct {
		Cursor     uint64
		State      uint64    // active, canceled, closed
		LastActive time.Time // track activity to deterime backoff
	}

	type endpoints struct {
		sync.Mutex
		*sync.Cond
		entries []endpoint
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

		endpoints endpoints

		err error
	}

	make := func(age time.Duration, length int, capacity ...int) *buffer {
		cap, ecap := uint64(length), uint64(32)
		switch {
		case len(capacity) >= 2:
			cap, ecap = uint64(capacity[0]), uint64(capacity[1])
		case len(capacity) == 1:
			cap = uint64(capacity[0])
		}
		len := uint64(length)
		if cap < len {
			cap = len
		}
		cap = uint64(1) << uint(math.Ceil(math.Log2(float64(cap))))
		ch := &buffer{
			len:   len,
			cap:   cap,
			age:   age,
			items: make([]item, cap),
			mod:   cap - 1,
			end:   cap,
			endpoints: endpoints{
				entries: make([]endpoint, 0, ecap),
			},
		}
		ch.endpoints.Cond = sync.NewCond(&ch.endpoints.Mutex)
		return ch
	}
	ch := make(age, length, capacity...)

	accessEndpoints := func(access func([]endpoint)) bool {
		spun := false
		for !atomic.CompareAndSwapUint32(&ch.endpoints.access, idle, busy) {
			runtime.Gosched()
			spun = true
		}
		access(ch.endpoints.entries)
		atomic.StoreUint32(&ch.endpoints.access, idle)
		return spun
	}

	send := func(value interface{}) {
		for ch.commit == ch.end {
			slowest := parked
			spun := accessEndpoints(func(endpoints []endpoint) {
				for i := range endpoints {
					cursor := atomic.LoadUint64(&endpoints[i].Cursor)
					if cursor < slowest {
						slowest = cursor
					}
				}
				if atomic.LoadUint64(&ch.begin) < slowest && slowest <= atomic.LoadUint64(&ch.end) {
					atomic.StoreUint64(&ch.begin, slowest)
					atomic.StoreUint64(&ch.end, slowest+ch.mod+1)
				} else {
					slowest = parked
				}
			})
			if slowest == parked {
				// no endpoints present...
				if !spun {
					// if we did not spin to access endpoints, then spin now.
					runtime.Gosched()
				}
				if atomic.LoadUint64(&ch.state) != active {
					return // channel no longer active
				}
			}
		}
		ch.items[ch.commit&ch.mod] = item{Value: value, At: time.Now()}
		atomic.AddUint64(&ch.commit, 1)
		ch.endpoints.Broadcast()
	}

	close := func(err error) {
		if atomic.CompareAndSwapUint64(&ch.state, active, closing) {
			ch.err = err
			if atomic.CompareAndSwapUint64(&ch.state, closing, closed) {
				accessEndpoints(func(endpoints []endpoint) {
					for i := range endpoints {
						atomic.CompareAndSwapUint64(&endpoints[i].State, active, closed)
					}
				})
			}
		}
		ch.endpoints.Broadcast()
	}

	observer := func(next interface{}, err error, done bool) {
		if atomic.LoadUint64(&ch.state) == active {
			if !done {
				send(next)
			} else {
				close(err)
			}
		}
	}

	appendEndpoint := func(cursor uint64) (ep *endpoint, err error) {
		accessEndpoints(func([]endpoint) {
			e := &ch.endpoints
			if len(e.entries) < cap(e.entries) {
				e.entries = append(e.entries, endpoint{Cursor: cursor})
				ep = &e.entries[len(e.entries)-1]
				return
			}
			for i := range e.entries {
				ep = &e.entries[i]
				if atomic.CompareAndSwapUint64(&ep.Cursor, parked, cursor) {
					return
				}
			}
			err = ErrOutOfEndpoints
			return
		})
		return
	}

	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		cursor := atomic.LoadUint64(&ch.begin)
		ep, err := appendEndpoint(cursor)
		if err != nil {
			runner := subscribeOn.Schedule(func() {
				if subscriber.Subscribed() {
					observe(nil, err, true)
				}
			})
			subscriber.OnUnsubscribe(runner.Cancel)
			return
		}
		commit := atomic.LoadUint64(&ch.commit)
		begin := atomic.LoadUint64(&ch.begin)
		if begin+ch.len < commit {
			atomic.StoreUint64(&ep.Cursor, commit-ch.len)
		}
		atomic.StoreUint64(&ep.State, atomic.LoadUint64(&ch.state))
		ep.LastActive = time.Now()

		receiver := subscribeOn.ScheduleFutureRecursive(0, func(self func(time.Duration)) {
			commit := atomic.LoadUint64(&ch.commit)

			// wait for commit to move away from cursor position, indicating fresh data
			if ep.Cursor == commit {
				if atomic.CompareAndSwapUint64(&ep.State, canceled, canceled) {
					// endpoint canceled
					atomic.StoreUint64(&ep.Cursor, parked)
					return
				} else {
					// endpoint still active (not canceled)
					now := time.Now()
					if now.Before(ep.LastActive.Add(1 * ms)) {
						// spinlock for 1ms (in increments of 50us) when no data from sender is arriving
						self(50 * us) // 20kHz
						return
					} else if now.Before(ep.LastActive.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&ep.State, closed, closed) {
							// channel has been closed
							observe(nil, ch.err, true)
							atomic.StoreUint64(&ep.Cursor, parked)
							return
						}
						// spinlock between 1ms and 250ms (in increments of 500us) of no data from sender
						self(500 * us) // 2kHz
						return
					} else {
						if subscribeOn.IsConcurrent() {
							// Block goroutine on condition until notified.
							ch.endpoints.Lock()
							ch.endpoints.Wait()
							ch.endpoints.Unlock()
							ep.LastActive = time.Now()
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
			if atomic.LoadUint64(&ep.State) == canceled {
				atomic.StoreUint64(&ep.Cursor, parked)
				return
			}
			for ; ep.Cursor != commit; atomic.AddUint64(&ep.Cursor, 1) {
				item := &ch.items[ep.Cursor&ch.mod]
				if ch.age == 0 || item.At.IsZero() || time.Since(item.At) < ch.age {
					observe(item.Value, nil, false)
				}
				if atomic.LoadUint64(&ep.State) == canceled {
					atomic.StoreUint64(&ep.Cursor, parked)
					return
				}
			}

			// all caught up; record time and loop back to wait for fresh data
			ep.LastActive = time.Now()
			self(0)
		})
		subscriber.OnUnsubscribe(receiver.Cancel)

		subscriber.OnUnsubscribe(func() {
			atomic.CompareAndSwapUint64(&ep.State, active, canceled)
			ch.endpoints.Broadcast()
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
//jig:needs NewBuffer, Subject<Foo>, Observer As<Foo>Observer

// NewSubjectFoo creates a new Subject. After the subject is terminated, all
// subsequent subscriptions to the observable side will be terminated
// immediately with either an Error or Complete notification send to the
// subscribing client
// 
// Note that this implementation is blocking. When there are subscribers, the
// observable goroutine is blocked until all subscribers have processed the
// next, error or complete notification.
func NewSubjectFoo() SubjectFoo {
	observer, observable := NewBuffer(0, 0, 1, 16)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}

//jig:template MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:template NewReplaySubject<Foo>
//jig:needs NewBuffer, Subject<Foo>, Observer As<Foo>Observer, MaxReplayCapacity

// NewReplaySubjectFoo creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
func NewReplaySubjectFoo(bufferCapacity int, windowDuration time.Duration) SubjectFoo {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	observer, observable := NewBuffer(windowDuration, bufferCapacity)
	return SubjectFoo{observer.AsFooObserver(), observable.AsObservableFoo()}
}
