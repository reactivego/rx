package ToSlice

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	slice, err := Range(1, 9).ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// [1 2 3 4 5 6 7 8 9]
	// <nil>
}

func Example_error() {
	source := Throw(RxError("crash"))

	slice, err := source.ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// []
	// crash
}
