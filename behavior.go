package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// Behavior returns a multicasting Observable[T] for an underlying Observable[T]
// as a Connectable[T] type, like [Observable.Publish], but with replay-latest,
// conflating, non-backpressuring semantics suited to distributing current
// state:
//
//   - Replay-latest on subscribe. A new Subscriber immediately observes the
//     current value — the construction seed until the source emits, then the
//     most recent source value.
//   - Conflation for live Subscribers. A Subscriber that falls behind skips
//     intermediate values and converges on the latest; it is never stranded on
//     a stale value (as a dropping FIFO buffer would be) and never blocks the
//     source (as a blocking FIFO buffer would, the way [Multicast] does).
//   - No backpressure to the source. The source is never blocked by a slow
//     Subscriber, so the source's own nature governs what Subscribers see: a
//     fast cold source (e.g. From over many values) on a serial scheduler may
//     run to completion before any Subscriber is scheduled, so Subscribers
//     observe only the final value and completion; a source that emits over
//     time lets Subscribers track it, modulo conflation when one lags.
//
// Composes as a Connectable, so the connection lifecycle is managed by the
// standard combinators:
//
//	source.Behavior(seed).RefCount()                 // count-free; replays latest
//	source.Behavior(seed).AutoConnect(n)             // connect after n subscribers
//	source.StartWith(seed).Behavior(seed).RefCount() // StartWith composes too
//
// Unlike [Observable.Publish], whose [Multicast] core does not replay,
// Behavior's core delivers the current value to a late Subscriber, which is
// what makes [Connectable.RefCount] over Behavior count-free: a Subscriber that
// attaches after Connect still receives the latest value rather than nothing.
//
// Behavior is scheduler-orthogonal. The conflating value cell and terminal
// latch are lock-light and shared; only the per-Subscriber wait differs by
// scheduler kind — a concurrent Subscriber blocks on a wake signal on its own
// goroutine, while a serial (trampoline) Subscriber polls by rescheduling,
// since a serial scheduler has no park/unpark and a blocking receiver would
// freeze the single goroutine it shares with the source.
func (observable Observable[T]) Behavior(seed T) Connectable[T] {
	state := newBehaviorState(seed)
	connector := func(scheduler Scheduler, subscriber Subscriber) {
		observe := func(next T, err error, done bool) {
			if done {
				state.complete(err)
				subscriber.Unsubscribe()
			} else {
				state.next(next)
			}
		}
		observable(observe, scheduler, subscriber)
	}
	return Connectable[T]{Observable: state.subscribe, Connector: connector}
}

// behaviorIdlePoll is the reschedule interval a serial (trampoline) Subscriber
// uses while no new value is available. A Subscriber that just delivered a
// value reschedules immediately (0) to catch up to the latest, then falls back
// to this interval when idle. Concurrent Subscribers do not poll — they block
// on a wake signal — so this affects serial schedulers only.
const behaviorIdlePoll = 5 * time.Millisecond

// behaviorValue is an immutable snapshot of the current value and its sequence
// number. The sequence increases on every source emission (values and the
// terminal), so a Subscriber can tell "is this newer than what I last
// delivered?" with a single comparison, and conflation falls out: when the
// sequence jumps, the Subscriber delivers only the newest snapshot.
type behaviorValue[T any] struct {
	value T
	seq   uint64
}

// behaviorTerm is the latched terminal (error or completion) and its sequence
// number. It is kept separate from the current value so a lagging Subscriber
// still delivers the final value (then the terminal) rather than conflating the
// last value away under the terminal.
type behaviorTerm struct {
	err error
	seq uint64
}

// behaviorReceiver is the per-Subscriber state. lastSeen is touched only by that
// Subscriber's own delivery loop/task (a single accessor), so it needs no lock.
// wake is a coalescing, capacity-1 signal used by concurrent Subscribers; quit
// is closed on unsubscribe to release a blocked concurrent Subscriber.
type behaviorReceiver struct {
	lastSeen uint64
	wake     chan struct{}
	quit     chan struct{}
}

// behaviorState is the shared conflating multicast core behind
// [Observable.Behavior]. The current value and terminal are stored as atomic
// snapshots so receivers read them without taking the lock; the lock guards
// only sequence increment and the receiver set.
type behaviorState[T any] struct {
	mu   sync.Mutex
	seq  uint64
	cur  atomic.Pointer[behaviorValue[T]]
	term atomic.Pointer[behaviorTerm]
	done bool
	subs map[*behaviorReceiver]struct{}
}

func newBehaviorState[T any](seed T) *behaviorState[T] {
	b := &behaviorState[T]{subs: make(map[*behaviorReceiver]struct{})}
	b.seq = 1
	b.cur.Store(&behaviorValue[T]{value: seed, seq: 1})
	return b
}

// next stores value as the new current snapshot (overwriting any undelivered
// previous value — this is the conflation) and wakes every receiver. It never
// blocks, so it applies no backpressure to the source.
func (b *behaviorState[T]) next(value T) {
	b.mu.Lock()
	if b.done {
		b.mu.Unlock()
		return
	}
	b.seq++
	b.cur.Store(&behaviorValue[T]{value: value, seq: b.seq})
	receivers := b.snapshotLocked()
	b.mu.Unlock()
	for _, r := range receivers {
		wakeReceiver(r)
	}
}

// complete latches the terminal and wakes every receiver. Idempotent: only the
// first terminal takes effect.
func (b *behaviorState[T]) complete(err error) {
	b.mu.Lock()
	if b.done {
		b.mu.Unlock()
		return
	}
	b.seq++
	b.done = true
	b.term.Store(&behaviorTerm{err: err, seq: b.seq})
	receivers := b.snapshotLocked()
	b.mu.Unlock()
	for _, r := range receivers {
		wakeReceiver(r)
	}
}

// snapshotLocked copies the receiver set so wakes happen outside the lock. The
// caller must hold b.mu.
func (b *behaviorState[T]) snapshotLocked() []*behaviorReceiver {
	receivers := make([]*behaviorReceiver, 0, len(b.subs))
	for r := range b.subs {
		receivers = append(receivers, r)
	}
	return receivers
}

func (b *behaviorState[T]) add(r *behaviorReceiver) {
	b.mu.Lock()
	b.subs[r] = struct{}{}
	b.mu.Unlock()
}

func (b *behaviorState[T]) remove(r *behaviorReceiver) {
	b.mu.Lock()
	delete(b.subs, r)
	b.mu.Unlock()
}

// wakeReceiver delivers a coalescing, non-blocking signal: if a wake is already
// pending the new one is dropped, because the receiver re-reads the latest
// snapshot on wake anyway. Never blocks the sender.
func wakeReceiver(r *behaviorReceiver) {
	select {
	case r.wake <- struct{}{}:
	default:
	}
}

// subscribe is the Connectable's multicasting Observable[T]. Each Subscriber
// gets a receiver that replays the current value, then tracks subsequent values
// with conflation, then the terminal. The wait strategy forks on the
// Subscriber's scheduler kind.
func (b *behaviorState[T]) subscribe(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
	r := &behaviorReceiver{
		wake: make(chan struct{}, 1),
		quit: make(chan struct{}),
	}
	b.add(r)

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			b.remove(r)
			close(r.quit)
		})
	}

	// deliver emits at most the latest value and then, if latched, the terminal,
	// advancing lastSeen past whatever it delivers. Conflation: it delivers the
	// current snapshot (the newest), not each intermediate value. Returns
	// whether it made progress and whether it delivered the terminal.
	deliver := func() (progressed, terminated bool) {
		if v := b.cur.Load(); v != nil && v.seq > r.lastSeen {
			r.lastSeen = v.seq
			observe(v.value, nil, false)
			progressed = true
		}
		if t := b.term.Load(); t != nil && t.seq > r.lastSeen {
			r.lastSeen = t.seq
			var zero T
			observe(zero, t.err, true)
			return true, true
		}
		return progressed, false
	}

	if scheduler.IsConcurrent() {
		// Concurrent: a dedicated goroutine blocks on the wake signal between
		// deliveries. Blocking is safe — it stalls only this Subscriber.
		runner := scheduler.Schedule(func() {
			for {
				if !subscriber.Subscribed() {
					cleanup()
					return
				}
				if _, terminated := deliver(); terminated {
					cleanup()
					return
				}
				select {
				case <-r.wake:
				case <-r.quit:
					return
				}
			}
		})
		subscriber.OnUnsubscribe(func() {
			cleanup()
			runner.Cancel()
		})
		return
	}

	// Serial (trampoline): never block the shared goroutine. Reschedule
	// immediately after delivering (to catch up to the latest), or after the
	// idle poll interval when there is nothing new.
	runner := scheduler.ScheduleFutureRecursive(0, func(again func(time.Duration)) {
		if !subscriber.Subscribed() {
			cleanup()
			return
		}
		progressed, terminated := deliver()
		if terminated {
			cleanup()
			return
		}
		if progressed {
			again(0)
		} else {
			again(behaviorIdlePoll)
		}
	})
	subscriber.OnUnsubscribe(func() {
		cleanup()
		runner.Cancel()
	})
}
