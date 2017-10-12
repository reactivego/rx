package Defer

import "fmt"

func Example_defer() {

	count := 0
	source := DeferInt(func() ObservableInt {
		count++
		return FromInt(count)
	})

	source.SubscribeNext(func(next int) {
		fmt.Printf("observable %d\n", next)
	})

	source.SubscribeNext(func(next int) {
		fmt.Printf("observable %d\n", next)
	})

	// Output:
	// observable 1
	// observable 2
}
