package Subscribe

import (
	"errors"
	"fmt"

	"github.com/reactivego/subscriber"
)

// Shows how to directly subscribe to an observable by calling it as a function.
// Don't do this in production code, use the Subscribe method.
func Example_subscribeDirect() {
	// Create an observable of int values
	observable := CreateInt(func(observer IntObserver) {
		observer.Next(456)
		observer.Complete()
		// Error will be delivered because we are not using the Subscribe method
		// to subscribe to the observable.
		observer.Error(errors.New("error"))
	})

	// Observe function. For both error and completed, done is true.
	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	// Non standard way of subscribing to an observable. But it does illustrate
	// that an observable is just a function that can be called directly.
	observable(observe, CurrentGoroutineScheduler(), subscriber.New())

	// Note that this incorrectly delivers the error after complete. So don't
	// subscribe like this, use the Subscribe method.

	// Output:
	// 456
	// complete
	// error
}

// Shows how to correctly subscribe to an observable by calling the Subscribe
// method.
func Example_subscribe() {
	// Create an observable of int values
	observable := CreateInt(func(observer IntObserver) {
		observer.Next(456)
		observer.Complete()
		// Error will not be delivered by Subscribe, because when Subscribe got
		// the Complete it immediately canceled the subscription.
		observer.Error(errors.New("error"))
	})

	// Observe function. For both error and completed, done is true.
	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}

	}

	// Subscribe (by default) only returns when observable has completed.
	observable.Subscribe(observe)

	// Output:
	// 456
	// complete
}
