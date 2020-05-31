package CombineLatestMapTo

import (
	_ "github.com/reactivego/rx/generic"
)

func Example_combineLatestMapToString() {
	From(1, 2, 3, 4).CombineLatestMapToString(OfString("Hey", "Ho")).Println()

	// Output:
	// [Hey Hey Hey Hey]
	// [Ho Hey Hey Hey]
	// [Ho Ho Hey Hey]
	// [Ho Ho Ho Hey]
	// [Ho Ho Ho Ho]
}
