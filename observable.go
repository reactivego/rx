package x

type Observable[T any] func(Observer[T], Scheduler, Subscriber)
