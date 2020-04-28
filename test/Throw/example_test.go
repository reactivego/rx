package Throw

import "fmt"

func Example_basic() {
	err := ThrowInt(RxError("throw")).Println()

	fmt.Println(err)
	// Output:
	// throw
}
