package CombineLatestWith

import _ "github.com/reactivego/rx/generic"

func Example_combineLatestWith() {

	A := From(1, 4, 7)
	B := From("two", "five", "eight")
	C := From(3, 6, 9)

	A.CombineLatestWith(B, C).Println()

	// Output:
	// [1 two 3]
	// [4 two 3]
	// [4 five 3]
	// [4 five 6]
	// [7 five 6]
	// [7 eight 6]
	// [7 eight 9]
}
