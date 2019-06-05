package Range

import (
	"fmt"
)

func Example_range() {
	Range(1, 12).
		Filter(func(x int) bool { return x%2 == 1 }).
		MapInt(func(x int) int { return x + x }).
		Println()

	//Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// 22
}

func Example_subscribe() {
	Range(1, 5).Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			fmt.Println("done")
		}
	})

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// done
}
