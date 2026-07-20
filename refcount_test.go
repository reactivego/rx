package rx

import (
	"context"
	"testing"
	"time"
)

// TestRefCountSurvivesFirstSubscriberUnsubscribe: the shared connection must
// stay up while any subscriber remains, even after the subscriber that
// triggered the connect (the first one) unsubscribes. The connection is torn
// down only when the refcount drops to zero.
func TestRefCountSurvivesFirstSubscriberUnsubscribe(t *testing.T) {
	feed := make(chan int)
	source := Create(func(index int) (int, error, bool) {
		next, ok := <-feed
		if !ok {
			return 0, nil, true
		}
		return next, nil, false
	})
	shared := source.Behavior(-1).RefCount()

	first := make(chan int, 8)
	sub1 := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			select {
			case first <- v:
			default:
			}
		}
	})

	second := make(chan int, 8)
	sub2 := shared.Subscribe(SchedulerContextWith(context.Background(), Goroutine), func(v int, err error, done bool) {
		if !done && err == nil {
			select {
			case second <- v:
			default:
			}
		}
	})
	defer sub2.Unsubscribe()

	feed <- 1
	drainUntil(t, first, 1, time.Second)
	drainUntil(t, second, 1, time.Second)

	// First subscriber leaves; refcount is still 1, so the shared connection
	// must keep pulling the source.
	sub1.Unsubscribe()

	select {
	case feed <- 2:
	case <-time.After(time.Second):
		t.Fatal("source no longer pulled: shared connection was torn down when the first subscriber left")
	}
	drainUntil(t, second, 2, time.Second)
}
