package Create

import (
	"errors"
	"fmt"

	"github.com/reactivego/subscriber"
)

// Shows how to directly subscribe to an observable by calling it as a function.
// Don't do this in a real application, use the Subscribe method.
func Example_directSubscribe() {
	// Create an observable of int values
	observable := CreateInt(func(observer IntObserver) {
		observer.Next(456)
		observer.Complete()
		// Error will be delivered because we are not using the Subscribe method
		// to subscribe to the observable.
		observer.Error(errors.New("error"))
	})

	observe := func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("complete")
			}
		}
	}

	// Non standard way of subscribing to an observable. But it shows that an
	// observable is just a function that can be called directly.
	observable(observe, NewTrampoline(), subscriber.New())

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

	observe := func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("complete")
			}
		}
	}
	observable.Subscribe(observe)

	// Output:
	// 456
	// complete
}
