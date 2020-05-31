package CatchError

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_catchError() {
	o123 := FromInt(1, 2, 3)
	o45 := FromInt(4, 5)
	oThrowError := ThrowInt(RxError("error"))
	err := o123.ConcatWith(oThrowError).CatchError(func(err error, caught ObservableInt) ObservableInt {
		if err == RxError("error") {
			return o45
		} else {
			return caught
		}
	}).Println()
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
