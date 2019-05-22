package Defer

import "fmt"

// Defer is used to postpone creating an observable until the moment somebody
// actually subscribes.
func Example_defer() {

	count := 0
	source := DeferInt(func() ObservableInt {
		count++
		return FromInt(count)
	})

	observe := func(next int) {
		fmt.Printf("observable %d\n", next)
	}
	source.SubscribeNext(observe)
	source.SubscribeNext(observe)

	// Output:
	// observable 1
	// observable 2
}
