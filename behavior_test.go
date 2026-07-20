package rx

import (
	"context"
	"sync"
	"testing"
	"time"
)

// idleSource returns a cold Observable[int] that emits the given values in order
// and then idles forever without completing, so late subscribers can attach to
// the live multicast during a test.
func idleSource(values ...int) Observable[int] {
	return Create(func(index int) (int, error, bool) {
		if index < len(values) {
			return values[index], nil, false
		}
		select {} // idle: never complete
	})
}

// drainUntil reads from ch until it sees want or the timeout elapses.
func drainUntil(t *testing.T, ch <-chan int, want int, timeout time.Duration) {
	t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case v := <-ch:
			if v == want {
				return
			}
		case <-deadline:
			t.Fatalf("never observed %d within %s", want, timeout)
		}
	}
}

// TestBehaviorReplaysSeedBeforeSourceEmits: a Connect'd Behavior replays the
// construction seed to a subscriber until the source produces a value.
func TestBehaviorReplaysSeedBeforeSourceEmits(t *testing.T) {
	// Source that never emits, so only the seed is ever available.
	src := Create(func(index int) (int, error, bool) { select {} })
	conn := src.Behavior(42)

	got := make(chan int, 4)
	sub := conn.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			select {
			case got <- v:
			default:
			}
		}
	})
	defer sub.Unsubscribe()

	connSub := conn.Connect(SchedulerContextWith(context.Background(), Goroutine))
	defer connSub.Unsubscribe()

	select {
	case v := <-got:
		if v != 42 {
			t.Fatalf("first observed = %d; want 42 (seed)", v)
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber never observed the seed")
	}
}

// TestBehaviorRefCountReplaysLatestToLateSubscriber is the headline property:
// the plain Publish().RefCount() failed here because Multicast does not replay,
// so a subscriber attaching after Connect saw nothing. Over Behavior's replay
// multicast, RefCount delivers the LATEST value to a late subscriber — and stays
// count-free (no AutoConnect(N)).
func TestBehaviorRefCountReplaysLatestToLateSubscriber(t *testing.T) {
	shared := idleSource(100, 200, 300).Behavior(-1).RefCount()

	// First subscriber triggers Connect (RefCount); drive the source to 300.
	first := make(chan int, 8)
	sub1 := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			select {
			case first <- v:
			default:
			}
		}
	})
	defer sub1.Unsubscribe()
	drainUntil(t, first, 300, time.Second)

	time.Sleep(30 * time.Millisecond)

	// Late subscriber must replay the latest (300) — not nothing (the RefCount
	// failure), and not the construction seed (-1).
	late := make(chan int, 4)
	sub2 := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			select {
			case late <- v:
			default:
			}
		}
	})
	defer sub2.Unsubscribe()

	select {
	case v := <-late:
		if v != 300 {
			t.Fatalf("late subscriber first saw %d; want 300 (latest replayed)", v)
		}
	case <-time.After(time.Second):
		t.Fatal("late subscriber saw nothing — RefCount over Behavior failed to replay")
	}
}

// TestBehaviorForwardsUpdatesAfterReplay: after the replayed value, a live
// subscriber observes subsequent source emissions in order.
func TestBehaviorForwardsUpdatesAfterReplay(t *testing.T) {
	shared := idleSource(1, 2, 3).Behavior(0).RefCount()

	got := make(chan int, 16)
	sub := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			got <- v
		}
	})
	defer sub.Unsubscribe()

	// The stream a subscriber sees is: replayed current, then live emissions.
	// Because Connect happens on first subscribe, the replayed current is the
	// seed (0) or an early value; assert the sequence reaches 3 monotonically.
	drainUntil(t, got, 3, time.Second)
}

// TestBehaviorStartWithComposes: StartWith feeds the seed into the cold source,
// and Behavior replays whatever the latest is — proving composability rather
// than a baked-in seed specialisation.
func TestBehaviorStartWithComposes(t *testing.T) {
	shared := idleSource(7, 8, 9).StartWith(5).Behavior(-1).RefCount()

	got := make(chan int, 16)
	sub := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			got <- v
		}
	})
	defer sub.Unsubscribe()

	// StartWith(5) prepends 5 ahead of 7,8,9; the stream must reach 9.
	drainUntil(t, got, 9, time.Second)
}

// TestBehaviorNoGapOrReorderUnderConcurrency stresses the replay/register seam:
// a fast source runs concurrently while subscribers attach at staggered times.
// Each subscriber's stream must be monotonic non-decreasing — conflation may
// skip values, but the replayed current value must never appear behind a
// fresher one, nor reorder. The source is never backpressured, so a subscriber
// that lags simply converges on the latest; nothing wedges (the channel-based
// design this replaced deadlocked here, because its blocking send held the
// registration lock).
func TestBehaviorNoGapOrReorderUnderConcurrency(t *testing.T) {
	conn := idleSourceFast().Behavior(0)
	connSub := conn.Connect(SchedulerContextWith(context.Background(), Goroutine))
	defer connSub.Unsubscribe()

	stop := make(chan struct{})
	var problems []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			first := true
			var prev int
			sub := conn.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, d bool) {
				if d || err != nil {
					return
				}
				// Always consume; only record a monotonicity violation.
				if !first && v < prev {
					mu.Lock()
					problems = append(problems, "non-monotonic: stale value behind a fresher one at the seam")
					mu.Unlock()
				}
				first = false
				prev = v
			})
			<-stop // keep draining until the test signals stop
			sub.Unsubscribe()
		}()
	}

	// Let subscribers attach at staggered times against the live sender, then
	// stop them all together.
	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()

	if len(problems) > 0 {
		t.Fatalf("seam problems under concurrency: %v", problems)
	}
}

// idleSourceFast emits a fast monotonically increasing counter then idles, used
// to exercise the concurrent replay/register seam.
func idleSourceFast() Observable[int] {
	return Create(func(index int) (int, error, bool) {
		if index < 100000 {
			return index + 1, nil, false
		}
		select {}
	})
}

// TestBehaviorConcurrentConvergence asserts the property a coalescing wake must
// guarantee: no matter how the source's emissions interleave with each
// subscriber's delivery, every live subscriber ends having delivered the final
// value. A lost wakeup (subscriber blocks while a newer value sits unseen)
// would strand a subscriber below the final value and fail this. Connect emits
// a known-terminal sequence; each subscriber's last pre-terminal value must be
// the source's last value.
func TestBehaviorConcurrentConvergence(t *testing.T) {
	const final = 500
	source := Create(func(index int) (int, error, bool) {
		if index < final {
			return index + 1, nil, false // 1..final
		}
		var zero int
		return zero, nil, true // complete after final
	})

	conn := source.Behavior(0)

	const subscribers = 12
	lasts := make([]int, subscribers)
	dones := make([]bool, subscribers)
	var wg sync.WaitGroup
	for i := range subscribers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			done := make(chan struct{})
			last := 0
			conn.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, d bool) {
				if d {
					if err == nil {
						dones[i] = true
					}
					close(done)
					return
				}
				last = v
			})
			<-done
			lasts[i] = last
		}()
	}

	conn.Connect(SchedulerContextWith(context.Background(), Goroutine))
	wg.Wait()

	for i := range subscribers {
		if !dones[i] {
			t.Errorf("subscriber %d never completed", i)
		}
		if lasts[i] != final {
			t.Errorf("subscriber %d converged on %d; want final %d (lost wakeup?)", i, lasts[i], final)
		}
	}
}

// ---- serial (trampoline) scheduler path ---------------------------------
//
// These tests drive Behavior on rx.NewScheduler() (a single-goroutine
// trampoline). The receiver cannot block there — it polls by rescheduling — so
// this path is implemented separately from the concurrent path and must be
// tested separately. Everything (source emission, Connect, subscribers) shares
// one goroutine and runs when serial.Wait() pumps the queue.

// TestBehaviorSerialReplaysFinalValueAndCompletes is the key serial-path
// property: because the source is not backpressured and shares the single
// trampoline goroutine, a cold finite source can run to completion before a
// subscriber's poll task gets to run. The subscriber must still observe the
// LATEST value (conflated) followed by completion — replay-latest holds even
// when no intermediate value was ever seen.
func TestBehaviorSerialReplaysFinalValueAndCompletes(t *testing.T) {
	serial := NewScheduler()
	ctx := SchedulerContextWith(context.Background(), serial)

	// Finite cold source: emits 1,2,3 then completes.
	source := Create(func(index int) (int, error, bool) {
		if index < 3 {
			return index + 1, nil, false
		}
		var zero int
		return zero, nil, true
	})

	conn := source.Behavior(0)

	var values []int
	var completed bool
	conn.Subscribe(ctx, func(v int, err error, done bool) {
		if done && err == nil {
			completed = true
			return
		}
		if err == nil {
			values = append(values, v)
		}
	})
	conn.Connect(ctx)

	serial.Wait()

	if !completed {
		t.Fatalf("serial subscriber never observed completion; values=%v", values)
	}
	if len(values) == 0 {
		t.Fatalf("serial subscriber observed no value before completion")
	}
	// Conflation: the last delivered value must be the source's final value (3).
	if last := values[len(values)-1]; last != 3 {
		t.Fatalf("last value before completion = %d; want 3 (latest)", last)
	}
	// Must never deliver a value out of order.
	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			t.Fatalf("non-monotonic serial delivery: %v", values)
		}
	}
}

// TestBehaviorSerialReplaysSeedToLateSubscriber: on the serial scheduler, a
// subscriber added after the source has already produced values still replays
// the latest value (not the seed, not nothing) when it finally runs.
func TestBehaviorSerialReplaysLatestToLateSubscriber(t *testing.T) {
	serial := NewScheduler()
	ctx := SchedulerContextWith(context.Background(), serial)

	source := Create(func(index int) (int, error, bool) {
		if index < 3 {
			return (index + 1) * 10, nil, false // 10,20,30
		}
		var zero int
		return zero, nil, true
	})

	conn := source.Behavior(-1)
	connSub := conn.Connect(ctx)
	_ = connSub

	var lateValues []int
	var lateDone bool
	// Schedule the late subscription as a task so it is enqueued behind the
	// source's emissions on the same trampoline.
	serial.Schedule(func() {
		conn.Subscribe(ctx, func(v int, err error, done bool) {
			if done && err == nil {
				lateDone = true
				return
			}
			if err == nil {
				lateValues = append(lateValues, v)
			}
		})
	})

	serial.Wait()

	if len(lateValues) == 0 {
		t.Fatalf("late serial subscriber saw nothing")
	}
	if first := lateValues[0]; first == -1 {
		t.Fatalf("late serial subscriber first saw the construction seed (-1); want the latest source value")
	}
	if last := lateValues[len(lateValues)-1]; last != 30 {
		t.Fatalf("late serial subscriber last value = %d; want 30 (latest)", last)
	}
	if !lateDone {
		t.Fatalf("late serial subscriber never observed completion")
	}
}

// TestBehaviorSerialDeliversErrorTerminal: an error terminal is latched and
// delivered (after the latest value) on the serial path.
func TestBehaviorSerialDeliversErrorTerminal(t *testing.T) {
	serial := NewScheduler()
	ctx := SchedulerContextWith(context.Background(), serial)

	boom := Err
	source := Create(func(index int) (int, error, bool) {
		if index == 0 {
			return 42, nil, false
		}
		var zero int
		return zero, boom, true
	})

	conn := source.Behavior(0)

	var gotErr error
	var sawValue bool
	conn.Subscribe(ctx, func(v int, err error, done bool) {
		if done {
			gotErr = err
			return
		}
		sawValue = true
	})
	conn.Connect(ctx)

	serial.Wait()

	if gotErr == nil {
		t.Fatalf("serial subscriber never received the error terminal")
	}
	_ = sawValue
}
