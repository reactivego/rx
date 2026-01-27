package rx

import "context"

// Context keys for storing rx-specific values
type contextKey int

const (
	schedulerKey contextKey = iota
)

// BackgroundGoroutine holds a context based on context.Background() with the
// Goroutine scheduler attached.
var BackgroundGoroutine = Background(Goroutine)

// Background returns a new context based on context.Background() with
// the provided scheduler attached.
func Background(scheduler Scheduler) context.Context {
	return WithScheduler(context.Background(), scheduler)
}

// WithScheduler returns a new context with the provided scheduler overriding
// any currently attached ones.
func WithScheduler(ctx context.Context, scheduler Scheduler) context.Context {
	return context.WithValue(ctx, schedulerKey, scheduler)
}

// withScheduler returns a new context based on the provided context with
// the provided scheduler attached (but only if none was already attached).
// If the provided context already has a scheduler attached, it returns the
// context as is.
// If no context is provided, then context.Background() is used.
// If no scheduler is provided, then NewScheduler() is used.
func withScheduler(ctx []context.Context, scheduler ...Scheduler) context.Context {
	if len(ctx) != 0 && ctx[0] != nil {
		if _, ok := ctx[0].Value(schedulerKey).(Scheduler); ok {
			return ctx[0]
		}
	} else {
		ctx = []context.Context{context.Background()}
	}
	if len(scheduler) != 0 {
		return context.WithValue(ctx[0], schedulerKey, scheduler[0])
	} else {
		return context.WithValue(ctx[0], schedulerKey, NewScheduler())
	}
}

// schedulerFromContext retrieves the scheduler from the context.
// It panics if the context does not have a scheduler attached.
// This function should only be called after withScheduler has been used
// to ensure a scheduler is present in the context.
func schedulerFromContext(ctx context.Context) Scheduler {
	return ctx.Value(schedulerKey).(Scheduler)
}
