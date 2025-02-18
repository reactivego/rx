package Wait

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	// Create an observable of int values
	observable := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(456)
		C()
		// Error will not be delivered by Subscribe, because when Subscribe got
		// the Complete it immediately canceled the subscription.
		E(RxError("error"))
	})

	// Subscribe (by default) only returns when observable has completed.
	err := observable.Wait()

	fmt.Println(err)
	// Output:
	// <nil>
}
