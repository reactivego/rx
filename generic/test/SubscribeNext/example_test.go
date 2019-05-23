package SubscribeNext

import (
	"errors"
	"fmt"
)

// Shows how to correctly subscribe to an observable by calling the
// SubscribeNext method.
func Example_subscribeNext() {
	// Create an observable of int values
	observable := CreateInt(func(observer IntObserver) {
		observer.Next(456)
		observer.Complete()
		// Error will not be delivered by Subscribe, because when Subscribe got
		// the Complete it immediately canceled the subscription.
		observer.Error(errors.New("error"))
	})

	// ObserveNext function only called for next values not for error or complete.
	observeNext := func(next int) {
		fmt.Println(next)
	}

	// SubscribeNext (by default) only returns when observable has completed.
	observable.SubscribeNext(observeNext)

	// Output:
	// 456
}
