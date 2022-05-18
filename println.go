package x

import (
	"fmt"
	"io"
)

func (observable Observable[T]) Println(schedulers ...Scheduler) error {
	return observable.Subscribe(Println[T](), schedulers...).Wait()
}

func Println[T any]() Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			fmt.Println(next)
		}
	}
}

func Fprintln[T any](out io.Writer) Observer[T] {
	return func(next T, err error, done bool) {
		if !done {
			fmt.Fprintln(out, next)
		}
	}
}
