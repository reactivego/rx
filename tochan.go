package x

func ToChan[T any](ch chan<- T) Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			ch <- next
		} else {
			close(ch)
		}
	}
}
