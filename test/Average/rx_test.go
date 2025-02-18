package Average

import _ "github.com/reactivego/rx/generic"

func Example_averageInt() {
	FromInt(1, 2, 3, 4, 5).Average().Println()

	// Output:
	// 3
}

func Example_averageFloat32() {
	FromFloat32(1, 2, 3, 4).Average().Println()

	// Output:
	// 2.5
}
