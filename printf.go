package rx

import "fmt"

func Printf[T any](format string) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Printf(format, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Printf(format string) Observable[T] {
	return Printf[T](format)(observable)
}
