package rx

// Creator[T] is a function type that generates values for an Observable stream.
//
// The Creator function receives a zero-based index for the current iteration and
// returns a tuple containing:
//   - Next: The next value to emit (of type T)
//   - Err: Any error that occurred during value generation
//   - Done: Boolean flag indicating whether the sequence is complete
//
// When Done is true, the Observable will complete after emitting any provided error.
// When Err is non-nil, the Observable will emit the error and then complete.
type Creator[T any] func(index int) (Next T, Err error, Done bool)
