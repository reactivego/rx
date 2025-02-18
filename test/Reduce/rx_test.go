package Reduce

import _ "github.com/reactivego/rx/generic"

func Example_reduce() {
	add := func(acc float32, value int) float32 {
		return acc + float32(value)
	}
	FromInt(1, 2, 3, 4, 5).ReduceFloat32(add, 0.0).Println()
	// Output: 15
}

func Example_reduceWeak() {
	add := func(acc interface{}, value int) interface{} {
		accint, ok := acc.(int)
		if !ok {
			return acc
		}
		return accint + value
	}
	FromInt(1, 2, 3, 4, 5).Reduce(add, 0).Println()
	// Output: 15
}
