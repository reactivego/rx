package observable

func Assign[T any](value *T) Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			*value = next
		}
	}
}
