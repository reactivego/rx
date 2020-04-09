package ToSlice

import (
	"fmt"
)

func Example_once() {
	slice, err := Range(1, 9).ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// [1 2 3 4 5 6 7 8 9]
	// <nil>
}

func Example_twice() {
	source := FromSliceInt([]int{1, 2, 3, 4})

	result, err := source.ToSlice()
	fmt.Println(result)
	fmt.Println(err)

	slice, err = source.ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// [1 2 3 4]
	// <nil>
	// [1 2 3 4]
	// <nil>
}
