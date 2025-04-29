package rx

type Observer[T any] func(next T, err error, done bool)

// Ignore creates an Observer that simply discards any emissions
// from an Observable. It is useful when you need to create an
// Observer but don't care about its values.
func Ignore[T any]() Observer[T] {
	return func(next T, err error, done bool) {}
}

// OnNext creates an Observer that only calls the provided function when a new
// value is emitted. It ignores completion and error signals.
func OnNext[T any](onNext func(T)) Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			onNext(next)
		}
	}
}

// OnDone creates an Observer that only calls the provided function when the
// Observable completes or errors. It passes any error that occurred to the callback.
func OnDone[T any](onDone func(error)) Observer[T] {
	return func(next T, err error, done bool) {
		if done {
			onDone(err)
		}
	}
}

// OnError creates an Observer that only calls the provided function when an
// error occurs. It ignores normal values and completion signals.
func OnError[T any](onError func(error)) Observer[T] {
	return func(next T, err error, done bool) {
		if done && err != nil {
			onError(err)
		}
	}
}

// OnComplete creates an Observer that only calls the provided function when the
// Observable completes. It ignores normal values and error signals.
func OnComplete[T any](onComplete func()) Observer[T] {
	return func(next T, err error, done bool) {
		if done && err == nil {
			onComplete()
		}
	}
}

// Next sends a new value to the Observer.
// It indicates that a new value has been emitted by the Observable.
func (observe Observer[T]) Next(next T) {
	observe(next, nil, false)
}

// Done signals that the Observable has completed emitting values,
// optionally with an error. If err is nil, it indicates normal completion.
// If err is non-nil, it indicates that the Observable terminated with an error.
func (observe Observer[T]) Done(err error) {
	var zero T
	observe(zero, err, true)
}
