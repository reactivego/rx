package Empty

import "fmt"

func Example_empty() {
	err := Empty().Println()
	fmt.Println(err)
	// Output: <nil>
}
