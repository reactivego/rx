package CombineAll

import _ "github.com/reactivego/rx"

type Vector struct{ x, y int }

func Example_combineAll() {
	a := FromInt(1, 2, 3)
	b := FromInt(4, 5, 6)

	FromObservableInt(a, b).CombineAll().MapVector(func(next IntSlice) Vector {
		return Vector{next[0], next[1]}
	}).Println()

	// Output:
	// {1 4}
	// {2 4}
	// {2 5}
	// {3 5}
	// {3 6}
}
