package rx

import (
	"context"

	"github.com/reactivego/scheduler"
)

// Scheduler is the interface for scheduling tasks, either serially or concurrently.
type Scheduler = scheduler.Scheduler

// ConcurrentScheduler is the interface for schedulers that run tasks concurrently.
type ConcurrentScheduler = scheduler.ConcurrentScheduler

// Goroutine is a concurrent scheduler that runs every task on its own goroutine.
var Goroutine = scheduler.Goroutine

// NewScheduler returns a new serial scheduler that runs tasks sequentially on the calling goroutine.
var NewScheduler = scheduler.New

// SchedulerContextWith returns a context with the scheduler attached, replacing any scheduler already present.
var SchedulerContextWith = scheduler.ContextWith

// SchedulerFromContext returns the Scheduler attached to the context, or nil and false when none is attached.
var SchedulerFromContext = scheduler.FromContext

// GoroutineContext returns a context that is guaranteed to have a scheduler
// attached. If no context or a nil context is provided, then
// context.Background() is used. If the context already has a scheduler it is
// replaced with scheduler.Goroutine.
func GoroutineContext(ctx ...context.Context) context.Context {
	if len(ctx) == 0 || ctx[0] == nil {
		ctx = []context.Context{context.Background()}
	}
	return scheduler.ContextWith(ctx[0], scheduler.Goroutine)
}

// contextWithScheduler returns a context that is guaranteed to have a scheduler
// attached. If the provided context already has a scheduler attached, it
// returns the context as is. If no context is provided, then
// context.Background() is used. If no scheduler is provided, then
// NewScheduler() is used.
func contextWithScheduler(ctx []context.Context, scheduler ...Scheduler) context.Context {
	if len(ctx) != 0 && ctx[0] != nil {
		if _, ok := SchedulerFromContext(ctx[0]); ok {
			return ctx[0]
		}
	} else {
		ctx = []context.Context{context.Background()}
	}
	if len(scheduler) != 0 {
		return SchedulerContextWith(ctx[0], scheduler[0])
	} else {
		return SchedulerContextWith(ctx[0], NewScheduler())
	}
}
