package SkipLast


import _ "github.com/reactivego/rx"

func Example_basic() {
	FromInt(1, 2, 3, 4, 5).SkipLast(2).Println()
	// Output:
	// 1
	// 2
	// 3
}
