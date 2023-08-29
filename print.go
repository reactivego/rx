package x

import (
	"fmt"
)

func Print[T any]() Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Print(next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Print() Observable[T] {
	return Print[T]()(observable)
}
