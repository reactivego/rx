package Defer

import "fmt"

func Example_defer() {
	// Defer is used to postpone creating an observable until the moment somebody
	// actually subscribes.
	count := 0
	source := DeferInt(func() ObservableInt {
		count++
		return FromInt(count)
	})
	mapped := source.MapString(func(next int) string {
		return fmt.Sprintf("observable %d", next)
	})

	// Println will subscribe to the observable
	mapped.Println()
	mapped.Println()

	// Output:
	// observable 1
	// observable 2
}
