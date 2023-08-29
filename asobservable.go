package x

func AsObservable[T any](observable Observable[any]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observable(observe.AsObserver(), scheduler, subscriber)
	}
}

func (observable Observable[T]) AsObservable() Observable[any] {
	return func(observe Observer[any], scheduler Scheduler, subscriber Subscriber) {
		observable(AsObserver[T](observe), scheduler, subscriber)
	}
}
