package x

func Combine[T any](observables ...Observable[T]) Observable[[]T] {
	return CombineAll(From(observables...))
}

func CombinePair[T, U any](first Observable[T], second Observable[U]) Observable[Pair[T, U]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable())), func(next []any) Pair[T, U] {
		return Pair[T, U]{First: next[0].(T), Second: next[1].(U)}
	})
}

func CombineTriple[T, U, V any](first Observable[T], second Observable[U], third Observable[V]) Observable[Triple[T, U, V]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable())), func(next []any) Triple[T, U, V] {
		return Triple[T, U, V]{First: next[0].(T), Second: next[1].(U), Third: next[2].(V)}
	})
}

func CombineQuadruple[T, U, V, W any](first Observable[T], second Observable[U], third Observable[V], fourth Observable[W]) Observable[Quadruple[T, U, V, W]] {
	return Map(CombineAll(From(first.AsObservable(), second.AsObservable(), third.AsObservable(), fourth.AsObservable())), func(next []any) Quadruple[T, U, V, W] {
		return Quadruple[T, U, V, W]{First: next[0].(T), Second: next[1].(U), Third: next[2].(V), Fourth: next[3].(W)}
	})
}
