package x

type Observer[T any] func(next T, err error, done bool)

func (observe Observer[T]) Next(next T) {
	observe(next, nil, false)
}

func (observe Observer[T]) Error(err error) {
	var zero T
	observe(zero, err, true)
}

func (observe Observer[T]) Complete() {
	var zero T
	observe(zero, nil, true)
}
