package rx

func Concat[T any](observables ...Observable[T]) Observable[T] {
	if len(observables) == 0 {
		return Empty[T]()
	}
	return observables[0].ConcatWith(observables[1:]...)
}
