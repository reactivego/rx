package x

func Equal[T comparable]() func(T, T) bool {
	return func(a T, b T) bool { return a == b }
}
