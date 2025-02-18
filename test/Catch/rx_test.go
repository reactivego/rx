package Catch

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_catch() {
	const problem = RxError("problem")

	err := FromInt(1, 2, 3).ConcatWith(ThrowInt(problem)).Catch(FromInt(4, 5)).Println()

	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
