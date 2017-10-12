package Reduce

import "fmt"

func Example_reduce() {
	add := func(acc float32, value int) float32 {
		return acc + float32(value)
	}
	FromInts(1, 2, 3, 4, 5).ReduceFloat32(add, 0.0).SubscribeNext(func(next float32) { fmt.Println(next) })

	// Output:
	// 15
}

func Example_reduceWeak() {
	add := func(acc interface{}, value int) interface{} {
		accint, ok := acc.(int)
		if !ok {
			return acc
		}
		return accint + value
	}
	FromInt(1, 2, 3, 4, 5).Reduce(add, 0).SubscribeNext(func(next interface{}) { fmt.Println(next) })

	// Output:
	// 15
}
