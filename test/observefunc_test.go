package test

import (
	"errors"
	"fmt"
)

// An ObserveFunc has Next, Error and Complete methods defined. But it is just
// a callback function meant to be passed to an observable.
func Example_observeFunc() {
	var observe ObserveFunc
	observe = func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	observe.Next(123)
	observe.Error(errors.New("error"))
	observe.Complete()

	// Output:
	// 123
	// error
	// complete
}
