package Take

import _ "github.com/reactivego/rx"

func Example_basic() {
	source := FromInt(1, 2, 3, 4, 5)
	source.Take(2).Println()
	source.Take(3).Println()
	// Output:
	// 1
	// 2
	// 1
	// 2
	// 3
}
