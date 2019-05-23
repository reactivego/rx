package Wait

import (
	"errors"
	"fmt"
)

// Shows how to correctly wait for an observable to finish without doing
// anything with the data.
func Example_wait() {
	// Create an observable of int values
	observable := CreateInt(func(observer IntObserver) {
		observer.Next(456)
		observer.Complete()
		// Error will not be delivered by Subscribe, because when Subscribe got
		// the Complete it immediately canceled the subscription.
		observer.Error(errors.New("error"))
	})

	// Subscribe (by default) only returns when observable has completed.
	err := observable.Wait()

	fmt.Println(err)

	// Output:
	// <nil>
}
