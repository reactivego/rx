package rx

// Ignore creates an Observer that simply discards any emissions
// from an Observable. It is useful when you need to create an
// Observer but don't care about its values.
func Ignore[T any]() Observer[T] {
	return func(next T, err error, done bool) {}
}
