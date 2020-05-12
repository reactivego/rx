package CombineLatestAll

import _ "github.com/reactivego/rx/generic"

type Vector struct{ x, y int }

func Example_combineAll() {
	a := FromInt(1, 2, 3)
	b := FromInt(4, 5, 6)

	FromObservableInt(a, b).CombineLatestAll().MapVector(func(next []int) Vector {
		return Vector{next[0], next[1]}
	}).Println()

	// Output:
	// {1 4}
	// {2 4}
	// {2 5}
	// {3 5}
	// {3 6}
}
