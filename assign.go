package rx

// Assign creates a pipe that assigns every next value that is not done to the
// provided variable. The pipe will forward all events (next, err, done) to the
// next observer.
func Assign[T any](value *T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					*value = next
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

// Assign is a method version of the Assign function.
// It assigns every next value that is not done to the provided variable.
func (observable Observable[T]) Assign(value *T) Observable[T] {
	return Assign(value)(observable)
}
