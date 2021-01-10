package StartWith

import (
	_ "github.com/reactivego/rx/generic"
)

func Example_startWith() {
	FromInt(2, 3).StartWith(1).Println()
	// Output:
	// 1
	// 2
	// 3
}
