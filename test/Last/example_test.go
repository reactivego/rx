package Last

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	err := FromInt(1, 2, 3, 4).Last().Println()

	fmt.Println(err)
	// Output:
	// 4
	// <nil>
}
