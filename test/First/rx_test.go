package First

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	err := FromInt(1, 2, 3, 4).First().Println()

	fmt.Println(err)

	// Output:
	// 1
	// <nil>
}
