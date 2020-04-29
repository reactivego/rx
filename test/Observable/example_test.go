package Observable

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

// An ObserveFunc is just a callback function meant to be passed to an
// observable.
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

	// Next
	observe(123, nil, false)
	// Error
	observe(zero, RxError("error"), true)
	// Complete
	observe(zero, nil, true)
	// Output:
	// 123
	// error
	// complete
}

// An observable is just a function, that takes an observe function, a scheduler,
// and a subscriber.
func Example_observable() {
	var observable Observable
	observable = func(observe ObserveFunc, s Scheduler, u Subscriber) {
		// Next
		observe(123, nil, false)
		// Error
		observe(zero, RxError("error"), true)
		// Complete
		observe(zero, nil, true)
	}

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

	observable(observe, nil, nil)
	// Output:
	// 123
	// error
	// complete
}
