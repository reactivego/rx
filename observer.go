package rx

import "errors"

// ErrTypecastFailed is returned when a type conversion fails during observer operations,
// typically when using AsObserver() to convert between generic and typed observers.
var ErrTypecastFailed = errors.Join(Err, errors.New("typecast failed"))

// Observer[T] represents a consumer of values delivered by an Observable.
// It is implemented as a function that takes three parameters:
// - next: the next value emitted by the Observable
// - err: any error that occurred during emission (nil if no error)
// - done: a boolean indicating whether the Observable has completed
//
// Observers follow the reactive pattern by receiving a stream of events
// (values, errors, or completion signals) and reacting to them accordingly.
type Observer[T any] func(next T, err error, done bool)

// Ignore creates an Observer that simply discards any emissions
// from an Observable. It is useful when you need to create an
// Observer but don't care about its values.
func Ignore[T any]() Observer[T] {
	return func(next T, err error, done bool) {}
}

// AsObserver converts an Observer of type `any` to an Observer of a specific type T.
// This allows adapting a generic Observer to a more specific type context.
func AsObserver[T any](observe Observer[any]) Observer[T] {
	return func(next T, err error, done bool) {
		observe(next, err, done)
	}
}

// AsObserver converts a typed Observer[T] to a generic Observer[any]. It
// handles type conversion from 'any' back to T, and will emit an
// ErrTypecastFailed error when conversion fails.
func (observe Observer[T]) AsObserver() Observer[any] {
	return func(next any, err error, done bool) {
		if !done {
			if nextT, ok := next.(T); ok {
				observe(nextT, err, done)
			} else {
				var zero T
				observe(zero, ErrTypecastFailed, true)
			}
		} else {
			var zero T
			observe(zero, err, true)
		}
	}
}

// Next sends a new value to the Observer. This is a convenience method that
// handles the common case of emitting a new value without errors or completion
// signals.
func (observe Observer[T]) Next(next T) {
	observe(next, nil, false)
}

// Done signals that the Observable has completed emitting values,
// optionally with an error. If err is nil, it indicates normal completion.
// If err is non-nil, it indicates that the Observable terminated with an error.
//
// After Done is called, the Observable will not emit any more values,
// regardless of whether the completion was successful or due to an error.
func (observe Observer[T]) Done(err error) {
	var zero T
	observe(zero, err, true)
}
