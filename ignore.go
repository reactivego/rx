package rx

func Ignore[T any]() Observer[T] {
	return func(next T, err error, done bool) {}
}
