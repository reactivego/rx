package CatchError

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_catchError() {
	const problem = RxError("problem")

	catcher := func(err error, caught ObservableInt) ObservableInt {
		if err == problem {
			return FromInt(4, 5)
		} else {
			return caught
		}
	}

	err := FromInt(1, 2, 3).ConcatWith(ThrowInt(problem)).CatchError(catcher).Println()

	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
