package Scan

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	add := func(acc int, value int) int {
		return acc + value
	}
	err := FromInt(1, 2, 3, 4, 5).ScanInt(add, 0).Println()
	fmt.Println(err)
	// Output:
	// 1
	// 3
	// 6
	// 10
	// 15
	// <nil>
}
