package observable

func Collect[T any](slice *[]T) Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			*slice = append(*slice, next)
		}
	}
}
