package ToSingle

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	value, err := FromInt(3).ToSingle()

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 3
	// <nil>
}

// ToSingle will return an error when the observable completes without emitting
// a single value.
func Example_emptyError() {
	value, err := EmptyInt().ToSingle()

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 0
	// expected one value, got none
}

// ToSingle will return an error when the observable emits multiple values.
func Example_multipleError() {
	value, err := FromInt(19, 20).ToSingle()

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 0
	// expected one value, got multiple
}
