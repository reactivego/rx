package CombineLatest

import _ "github.com/reactivego/rx"

func Example_basic() {

	A := FromInt(1, 3, 5)
	B := FromInt(2, 4, 6)

	CombineLatestInt(A, B).Println()

	// Output:
	// [1 2]
	// [3 2]
	// [3 4]
	// [5 4]
	// [5 6]
}
