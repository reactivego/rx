package observable

type Observable[T any] func(Observer[T], Scheduler, Subscriber)
