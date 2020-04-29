package ConcatAll

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_concatAll() {
	Interval(time.Millisecond).Take(3).MapObservableInt(func(next int) ObservableInt {
		return Range(next, 2)
	}).ConcatAll().Println()

	// Output:
	// 0
	// 1
	// 1
	// 2
	// 2
	// 3
}
