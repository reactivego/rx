package rx

// Append creates a pipe that appends each emitted value to the provided slice.
// It passes each value through to the next observer after appending it.
// This allows collecting all emitted values in a slice while still forwarding them.
// Only values emitted before completion (done=false) are appended.
func Append[T any](slice *[]T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					*slice = append(*slice, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

// Append is a method variant of the Append function that appends each emitted value
// to the provided slice while forwarding all emissions to downstream operators.
// This is a convenience method that calls the standalone Append function.
func (observable Observable[T]) Append(slice *[]T) Observable[T] {
	return Append(slice)(observable)
}
