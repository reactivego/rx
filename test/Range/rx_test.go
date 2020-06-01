package Range

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	Range(1, 3).Println()
	//Output:
	// 1
	// 2
	// 3
}

func Example_complex() {
	RangeInt(1, 12).
		Filter(func(x int) bool { return x%2 == 1 }).
		MapInt(func(x int) int { return x + x }).
		DoOnComplete(func() { fmt.Println("done") }).
		Println()
	//Output:
	// 2
	// 6
	// 10
	// 14
	// 18
	// 22
	// done
}
