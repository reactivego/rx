package TakeLast

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	source := FromInt(1, 2, 3, 4, 5)

	err := source.TakeLast(2).Println()
	fmt.Println(err)

	err = source.TakeLast(3).Println()
	fmt.Println(err)
	// Output:
	// 4
	// 5
	// <nil>
	// 3
	// 4
	// 5
	// <nil>
}
