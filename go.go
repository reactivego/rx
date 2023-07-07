package x

func (observable Observable[T]) Go() {
	observable.Subscribe(Ignore[T](), Goroutine)
}
