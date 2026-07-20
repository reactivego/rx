package rx

import (
	"context"
	"testing"
	"time"
)

// TestWaitReusesSerialSchedulerContext: a context carrying one serial
// scheduler must be reusable for consecutive Wait calls; the second Wait used
// to hang forever because the trampoline's task loop only ran once.
func TestWaitReusesSerialSchedulerContext(t *testing.T) {
	ctx := SchedulerContextWith(context.Background(), NewScheduler())
	if err := From(1, 2, 3).Wait(ctx); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- From(4, 5, 6).Wait(ctx) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("second Wait on the same serial scheduler context hung")
	}
}
