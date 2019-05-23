package Observable

import (
	"errors"
	"fmt"
)

// An observable is just a function, that takes an observe function, a scheduler,
// and a subscriber.
func Example_observable() {
	observable := func(observe ObserveFunc, s Scheduler, u Subscriber) {
		observe.Next(123)
		observe.Error(errors.New("error"))
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
