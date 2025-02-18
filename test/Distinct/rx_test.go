package Distinct

import _ "github.com/reactivego/rx/generic"

func Example_distinct() {
	FromInt(1, 2, 2, 1, 3).Distinct().Println()
	// Output:
	// 1
	// 2
	// 3
}
