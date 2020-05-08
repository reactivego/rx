package Create

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)
// Shows how to use CreateString to create an observable of strings
func Example_createString() {
	source := CreateString(func(N NextString, E Error, C Complete, X Canceled) {
		if X() {
			return
		}
		N("Hello")
		N("World!")
		C()
	})

	err := source.Println()

	if err == nil {
		fmt.Println("success")
	}

	// Output:
	// Hello
	// World!
	// success
}
