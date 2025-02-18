package Skip

import _ "github.com/reactivego/rx/generic"

func Example_basic() {
	FromInt(1, 2, 3, 4, 5).Skip(2).Println()
	// Output:
	// 3
	// 4
	// 5
}
