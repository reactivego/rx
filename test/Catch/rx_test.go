package Catch

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_catch() {
	throw := ThrowInt(RxError("error"))
	err := FromInt(1, 2, 3).ConcatWith(throw).Catch(FromInt(4, 5)).Println()
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
