package Map

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_mapString() {
	FromInt(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).Println()

	//Output:
	// 1!
	// 2!
	// 3!
	// 4!
}

func Example_mapVector() {
	//jig:type Vector []int

	FromVector([]int{1, 2, 3}, []int{5, 6, 7}).MapInt(func(x []int) int {
		sum := 0
		for _, v := range x {
			sum += v
		}
		return sum
	}).Println()

	//Output:
	// 6
	// 18
}
