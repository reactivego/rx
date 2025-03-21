package rx

func Defer[T any](factory func() Observable[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
}
