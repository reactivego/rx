package Min

import _ "github.com/reactivego/rx/generic"

func Example_basic() {
	FromInt(4, 5, 4, 3, 2, 1, 2).Min().Println()
	// Output: 1
}
