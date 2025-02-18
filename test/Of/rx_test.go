package Of

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	err := Of(1, 2, 3, 4, 5).Println()

	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
