package TakeWhile

import _ "github.com/reactivego/rx"

type predicate func(int) bool

func smallerThan(value int) predicate {
	return func(next int) bool {
		return next < value
	}
}

func Example_takeWhile() {
	FromInt(1, 2, 3, 4, 5).TakeWhile(smallerThan(4)).Println()
	// Output:
	// 1
	// 2
	// 3
}
