package rx

func Never[T any]() Observable[T] {
	return func(Observer[T], Scheduler, Subscriber) {}
}
