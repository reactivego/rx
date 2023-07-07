package x

import (
	"fmt"
	"io"
)

func Fprint[T any](out io.Writer) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					fmt.Fprint(out, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Fprint(out io.Writer) Observable[T] {
	return observable.Pipe(Fprint[T](out))
}
