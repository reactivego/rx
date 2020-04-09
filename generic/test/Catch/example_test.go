package Catch

import "fmt"

func Example_catch() {
	o123 := FromInts(1, 2, 3)
	o45 := FromInts(4, 5)
	oThrowError := ThrowInt(RxError("error"))

	err := o123.Concat(oThrowError).Catch(o45).Println()
	fmt.Println(err)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// <nil>
}
