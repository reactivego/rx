package rx

func MergeMap[T, U any](observable Observable[T], project func(T) Observable[U]) Observable[U] {
	return MergeAll(Map(observable, project))
}
