package rx

func EndWith[T any](values ...T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return observable.ConcatWith(From(values...))
	}
}

func (observable Observable[T]) EndWith(values ...T) Observable[T] {
	return EndWith(values...)(observable)
}
