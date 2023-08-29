package x

import (
	"fmt"
	"io"
)

func Fprintf[T any](out io.Writer, format string) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Fprintf(out, format, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Fprintf(out io.Writer, format string) Observable[T] {
	return Fprintf[T](out, format)(observable)
}
