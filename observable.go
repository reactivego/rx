package rx

type Observable[T any] func(Observer[T], Scheduler, Subscriber)
