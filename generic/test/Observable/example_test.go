package Observable

import (
	"fmt"
)

// An observable is just a function, that takes an observe function, a scheduler,
// and a subscriber.
func Example_observable() {
	observable := func(observe ObserveFunc, s Scheduler, u Subscriber) {
		observe.Next(123)
		observe.Error(RxError("error"))
		observe.Complete()
	}

	observe := func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	observable(observe, nil, nil)

	// Output:
	// 123
	// error
	// complete
}


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
	observe.Error(RxError("error"))
	observe.Complete()

	// Output:
	// 123
	// error
	// complete
}
