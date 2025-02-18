package Max

import _ "github.com/reactivego/rx/generic"

func Example_test() {
	FromInt(4, 5, 4, 3, 2, 1, 2).Max().Println()

	// Output: 5
}
