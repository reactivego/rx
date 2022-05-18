package x

func (observable Observable[T]) StartWith(values ...T) Observable[T] {
	return From(values...).ConcatWith(observable)
}
