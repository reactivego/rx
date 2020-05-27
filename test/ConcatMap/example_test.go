package ConcatMap

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_concatMapString() {
	FromInt(0, 1, 2, 3).ConcatMapString(func(v int) ObservableString {
		return JustString(fmt.Sprintf("v = %v", v))
	}).Println()

	//Output:
	// v = 0
	// v = 1
	// v = 2
	// v = 3
}
