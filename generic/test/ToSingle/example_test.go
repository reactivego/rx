package ToSingle

import (
	"fmt"
)

// ToSingle is used to make sure only a single value was produced by the
// observable. ToSingle will internally run the observable on an asynchronous
// scheduler. However it will only return when the observable was complete or
// an error was emitted.
func Example_toSingle() {
	if value, err := FromInt(19).ToSingle(); err == nil {
		fmt.Println(value)
	}

	if _, err := FromInts(19, 20).ToSingle(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 19
	// expected one value, got multiple
}
