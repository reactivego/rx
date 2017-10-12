package As

import (
	"fmt"
)

func Example_as() {
	// We are using From(...interface{}) to create an Observable of any value.

	// We should be able to convert [0.1, 1.2, 2.3] to float64 values using AsFloat64.
	From(0.1, 1.2, 2.3).AsFloat64().SubscribeNext(func(next float64) {
		fmt.Println(next)
	})

	fmt.Println("---")

	// We should not be able to convert "Hello, Rx!" to float64 values using AsFloat64.
	err := From("Hello, Rx!").AsFloat64().Wait()
	fmt.Println(err)

	// Output:
	// 0.1
	// 1.2
	// 2.3
	// ---
	// typecast to float64 failed
}

func Example_asAny() {
	// We convert an ObservableString to Observable.
	FromString("Hello, Rx!").AsAny().SubscribeNext(func(next interface{}) {
		fmt.Println(next)
	})

	// Output:
	// Hello, Rx!
}
