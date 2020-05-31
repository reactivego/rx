package Println

import _ "github.com/reactivego/rx"

// Subscribe to an observable using Println and print all emitted values.
func Example_println() {
	Range(0, 3).Println()
	// Output:
	// 0
	// 1
	// 2
}
