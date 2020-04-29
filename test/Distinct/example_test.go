package Distinct

import _ "github.com/reactivego/rx"

func Example_distinct() {
	FromInt(1, 1, 2, 2, 3, 2, 4, 5).Distinct().Println()
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
