package x

func ConcatMap[T, U any](observable Observable[T], project func(T) Observable[U]) Observable[U] {
	return ConcatAll(Map(observable, project))
}
