package Catch

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_catch() {
	o123 := FromInt(1, 2, 3)
	o45 := FromInt(4, 5)
	oThrowError := ThrowInt(RxError("error"))

	err := o123.ConcatWith(oThrowError).Catch(o45).Println()
	fmt.Println(err)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
