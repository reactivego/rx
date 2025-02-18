package Observable

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

// An Observer is just a callback function meant to be passed to an Observable.
func Example_observer() {
	var observe Observer
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
	observe(nil, RxError("error"), true)
	// Complete
	observe(nil, nil, true)
	// Output:
	// 123
	// error
	// complete
}

// An Observable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
func Example_observable() {
	var observable Observable
	observable = func(observe Observer, s Scheduler, u Subscriber) {
		// Next
		observe(123, nil, false)
		// Error
		observe(nil, RxError("error"), true)
		// Complete
		observe(nil, nil, true)
	}

	var observer Observer
	observer = func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	observable(observer, nil, nil)
	// Output:
	// 123
	// error
	// complete
}
