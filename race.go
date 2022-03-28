package observable

func Race[T any](observables ...Observable[T]) Observable[T] {
	if len(observables) == 0 {
		return Empty[T]()
	}
	return observables[0].RaceWith(observables[1:]...)
}
