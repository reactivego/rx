package ToSlice

import (
	"fmt"
)

// ToSlice will run the observable to completion and then return the resulting
// slice and any error that was emitted by the observable. ToSlice will
// internally run the observable on an asynchronous scheduler.
func Example_toSlice() {
	if slice, err := Range(1, 9).ToSlice(); err == nil {
		for _, value := range slice {
			fmt.Print(value)
		}
	}

	// Output:
	// 123456789
}
