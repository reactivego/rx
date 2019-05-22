package AsObservable

import (
	"fmt"
)

func Example_asObservableFloat64() {
	// We are using From(...interface{}) to create an Observable of any value.

	// We should be able to convert [0.1, 1.2, 2.3] to float64 values using AsFloat64.
	From(0.1, 1.2, 2.3).AsObservableFloat64().SubscribeNext(func(next float64) {
		fmt.Println(next)
	})

	fmt.Println("---")

	// We should not be able to convert "Hello, Rx!" to float64 values using AsObservableFloat64.
	err := From("Hello, Rx!").AsObservableFloat64().Wait()
	fmt.Println(err)

	// Output:
	// 0.1
	// 1.2
	// 2.3
	// ---
	// typecast to float64 failed
}

func Example_asObservable() {
	// We convert an ObservableString to Observable.
	FromString("Hello, Rx!").AsObservable().SubscribeNext(func(next interface{}) {
		fmt.Println(next)
	})

	// Output:
	// Hello, Rx!
}
