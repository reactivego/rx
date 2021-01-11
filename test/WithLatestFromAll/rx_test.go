package WithLatestFromAll

import _ "github.com/reactivego/rx/generic"

func Example_withLatestFromAll() {
	a := FromInt(1, 2, 3, 4, 5).AsObservable()
	b := FromString("A", "B", "C", "D", "E").AsObservable()
	c := FromObservable(a, b)
	c.WithLatestFromAll().Println()

	// Output:
	// [2 A]
	// [3 B]
	// [4 C]
	// [5 D]
}
