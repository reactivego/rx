package Throw

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	err := ThrowInt(RxError("throw")).Println()
	fmt.Println(err)
	// Output:
	// throw
}
