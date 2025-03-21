package rx

const TypecastFailed = Error("typecast failed")

func AsObserver[T any](observe Observer[any]) Observer[T] {
	return func(next T, err error, done bool) {
		observe(next, err, done)
	}
}

func (observe Observer[T]) AsObserver() Observer[any] {
	return func(next any, err error, done bool) {
		if !done {
			if nextT, ok := next.(T); ok {
				observe(nextT, err, done)
			} else {
				var zero T
				observe(zero, TypecastFailed, true)
			}
		} else {
			var zero T
			observe(zero, err, true)
		}
	}
}
