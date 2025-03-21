package rx

import "fmt"

func Println[T any]() Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Println(next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Println() Observable[T] {
	return Println[T]()(observable)
}
