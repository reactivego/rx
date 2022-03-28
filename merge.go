package observable

func Merge[T any](observables ...Observable[T]) Observable[T] {
	if len(observables) == 0 {
		return Empty[T]()
	}
	return observables[0].MergeWith(observables[1:]...)
}
