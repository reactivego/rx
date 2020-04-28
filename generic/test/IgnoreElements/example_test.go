package IgnoreElements

import
	"fmt"

func Example_ignoreElements() {
	err := FromInts(1, 2, 3, 4, 5).IgnoreElements().Println()
	fmt.Println("error", err)

	// Output:
	// error <nil>
}
