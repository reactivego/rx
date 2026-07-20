package rx

import (
	"sync/atomic"
	"testing"
	"time"
)

// countingSource returns a never-completing Observable[int] that counts how
// many times it is subscribed to, i.e. how often a Connectable connects.
func countingSource(connects *atomic.Int32) Observable[int] {
	src := Create(func(index int) (int, error, bool) { select {} })
	return func(observe Observer[int], scheduler Scheduler, subscriber Subscriber) {
		connects.Add(1)
		src(observe, scheduler, subscriber)
	}
}

// TestAutoConnectDoesNotReconnectWhileConnected: the refcount dipping below
// count and returning must not trigger a second connect while the first
// connection is still live.
func TestAutoConnectDoesNotReconnectWhileConnected(t *testing.T) {
	var connects atomic.Int32
	shared := countingSource(&connects).Behavior(-1).AutoConnect(2)

	subA := shared.Subscribe(GoroutineContext(), Ignore[int]())
	subB := shared.Subscribe(GoroutineContext(), Ignore[int]())
	defer subB.Unsubscribe()
	time.Sleep(50 * time.Millisecond)
	if got := connects.Load(); got != 1 {
		t.Fatalf("connects = %d after reaching count; want 1", got)
	}

	subA.Unsubscribe()
	subC := shared.Subscribe(GoroutineContext(), Ignore[int]())
	defer subC.Unsubscribe()
	time.Sleep(50 * time.Millisecond)
	if got := connects.Load(); got != 1 {
		t.Fatalf("connects = %d after refcount dipped below count and returned; want 1", got)
	}
}

// TestAutoConnectReconnectsAfterCancel: once every subscriber has left and the
// connection is canceled, a new wave of subscribers reaching count connects again.
func TestAutoConnectReconnectsAfterCancel(t *testing.T) {
	var connects atomic.Int32
	shared := countingSource(&connects).Behavior(-1).AutoConnect(2)

	subA := shared.Subscribe(GoroutineContext(), Ignore[int]())
	subB := shared.Subscribe(GoroutineContext(), Ignore[int]())
	time.Sleep(50 * time.Millisecond)
	subA.Unsubscribe()
	subB.Unsubscribe()
	time.Sleep(50 * time.Millisecond)

	subC := shared.Subscribe(GoroutineContext(), Ignore[int]())
	defer subC.Unsubscribe()
	subD := shared.Subscribe(GoroutineContext(), Ignore[int]())
	defer subD.Unsubscribe()
	time.Sleep(50 * time.Millisecond)
	if got := connects.Load(); got != 2 {
		t.Fatalf("connects = %d after cancel and a new wave; want 2", got)
	}
}
