package rx

func StartWith[T any](values ...T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return From(values...).ConcatWith(observable)
	}
}

func (observable Observable[T]) StartWith(values ...T) Observable[T] {
	return StartWith(values...)(observable)
}
