package rx

func Zip[T any](observables ...Observable[T]) Observable[[]T] {
	return ZipAll(From(observables...))
}

func Zip2[T, U any](first Observable[T], second Observable[U], options ...MaxBufferSizeOption) Observable[Tuple2[T, U]] {
	return Map(ZipAll(From(first.AsObservable(), second.AsObservable()), options...), func(next []any) Tuple2[T, U] {
		return Tuple2[T, U]{next[0].(T), next[1].(U)}
	})
}

func Zip3[T, U, V any](first Observable[T], second Observable[U], third Observable[V], options ...MaxBufferSizeOption) Observable[Tuple3[T, U, V]] {
	return Map(ZipAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable()), options...), func(next []any) Tuple3[T, U, V] {
		return Tuple3[T, U, V]{next[0].(T), next[1].(U), next[2].(V)}
	})
}

func Zip4[T, U, V, W any](first Observable[T], second Observable[U], third Observable[V], fourth Observable[W], options ...MaxBufferSizeOption) Observable[Tuple4[T, U, V, W]] {
	return Map(ZipAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable(), fourth.AsObservable()), options...), func(next []any) Tuple4[T, U, V, W] {
		return Tuple4[T, U, V, W]{next[0].(T), next[1].(U), next[2].(V), next[3].(W)}
	})
}

func Zip5[T, U, V, W, X any](first Observable[T], second Observable[U], third Observable[V], fourth Observable[W], fifth Observable[X], options ...MaxBufferSizeOption) Observable[Tuple5[T, U, V, W, X]] {
	return Map(ZipAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable(), fourth.AsObservable(), fifth.AsObservable()), options...), func(next []any) Tuple5[T, U, V, W, X] {
		return Tuple5[T, U, V, W, X]{next[0].(T), next[1].(U), next[2].(V), next[3].(W), next[4].(X)}
	})
}
