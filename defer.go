package rx

// Defer creates an Observable that will use the provided factory function to
// create a new Observable every time it's subscribed to. This is useful for
// creating cold Observables or for delaying expensive Observable creation until
// subscription time.
func Defer[T any](factory func() Observable[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
}
