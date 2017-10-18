package Create

import (
	"fmt"
)

// Shows how to use CreateString to create an observable of strings
func Example_createString() {
	source := CreateString(func(observer StringObserver) {
		observer.Next("Hello")
		observer.Next("World!")
		observer.Complete()
	})

	source.SubscribeNext(func(value string) { fmt.Print(value) })

	// Output:
	// HelloWorld!
}
