package x

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
