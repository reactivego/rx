package Error

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_toSlice() {
	slice, err := ErrorInt(RxError("rx-error")).ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// []
	// rx-error
}

func Example_println() {
	err := ErrorInt(RxError("rx-error")).Println()

	fmt.Println(err)
	// Output:
	// rx-error
}
