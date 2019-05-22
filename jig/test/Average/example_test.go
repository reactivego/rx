package Average

import "fmt"

func Example_averageInt() {
	FromInts(1, 2, 3, 4, 5).Average().SubscribeNext(func(next int) {
		fmt.Println(next)
	})

	// Output:
	// 3
}

func Example_averageFloat32() {
	FromFloat32s(1, 2, 3, 4).Average().SubscribeNext(func(next float32) {
		fmt.Println(next)
	})

	// Output:
	// 2.5
}
