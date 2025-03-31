package rx

func CombineLatest[T any](observables ...Observable[T]) Observable[[]T] {
	return CombineAll(From(observables...))
}

func CombineLatest2[T, U any](first Observable[T], second Observable[U]) Observable[Tuple2[T, U]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable())), func(next []any) Tuple2[T, U] {
		return Tuple2[T, U]{next[0].(T), next[1].(U)}
	})
}

func CombineLatest3[T, U, V any](first Observable[T], second Observable[U], third Observable[V]) Observable[Tuple3[T, U, V]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable())), func(next []any) Tuple3[T, U, V] {
		return Tuple3[T, U, V]{next[0].(T), next[1].(U), next[2].(V)}
	})
}

func CombineLatest4[T, U, V, W any](first Observable[T], second Observable[U], third Observable[V], fourth Observable[W]) Observable[Tuple4[T, U, V, W]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable(), fourth.AsObservable())), func(next []any) Tuple4[T, U, V, W] {
		return Tuple4[T, U, V, W]{next[0].(T), next[1].(U), next[2].(V), next[3].(W)}
	})
}

func CombineLatest5[T, U, V, W, X any](first Observable[T], second Observable[U], third Observable[V], fourth Observable[W], fifth Observable[X]) Observable[Tuple5[T, U, V, W, X]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable(), fourth.AsObservable(), fifth.AsObservable())), func(next []any) Tuple5[T, U, V, W, X] {
		return Tuple5[T, U, V, W, X]{next[0].(T), next[1].(U), next[2].(V), next[3].(W), next[4].(X)}
	})
}
