package ConcatMapTo

import _ "github.com/reactivego/rx"

func Example_basic() {
	FromInt(0, 1, 2, 3).ConcatMapToString(OfString("change")).Println()

	//Output:
	// change
	// change
	// change
	// change
}
