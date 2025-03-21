package rx

type Pipe[T any] func(Observable[T]) Observable[T]

func (observable Observable[T]) Pipe(segments ...Pipe[T]) Observable[T] {
	for _, s := range segments {
		observable = s(observable)
	}
	return observable
}
