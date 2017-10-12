package ConcatAll

import (
	"fmt"
)

func Example_concatAll() {
	source := Range(0, 3).MapObservableInt(func(next int) ObservableInt {
		return Range(next, 2)
	}).ConcatAll()

	source.SubscribeNext(func(next int) {
		fmt.Println(next)
	})

	// Output:
	// 0
	// 1
	// 1
	// 2
	// 2
	// 3
}
