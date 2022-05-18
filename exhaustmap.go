package x

func ExhaustMap[T, U any](observable Observable[T], project func(T) Observable[U]) Observable[U] {
	return ExhaustAll(Map(observable, project))
}
