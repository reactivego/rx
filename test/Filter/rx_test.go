package Filter

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	even := func(i int) bool {
		return i%2 == 0
	}

	err := FromInt(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).Println()
	fmt.Println(err)

	// Output:
	// 2
	// 4
	// 6
	// 8
	// <nil>
}
