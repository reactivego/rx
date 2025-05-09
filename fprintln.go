package rx

import (
	"fmt"
	"io"
)

func Fprintln[T any](out io.Writer) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Fprintln(out, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Fprintln(out io.Writer) Observable[T] {
	return Fprintln[T](out)(observable)
}
