package rx

import (
	"fmt"
)

func ExampleObservable_Map() {
	From(1, 2, 3, 4).Map(func(i any) any {
		return fmt.Sprintf("%d!", i.(int))
	}).Println()

	// Output:
	// 1!
	// 2!
	// 3!
	// 4!
}
