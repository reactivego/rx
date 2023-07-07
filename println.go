package x

import (
	"os"
)

func (observable Observable[T]) Println(schedulers ...Scheduler) error {
	return observable.Pipe(Fprintln[T](os.Stdout)).Wait(schedulers...)
}
