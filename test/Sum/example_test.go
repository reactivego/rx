package Sum

import _ "github.com/reactivego/rx"

func Example_sumInt() {
	FromInt(1, 2, 3, 4, 5).Sum().Println()
	// Output: 15
}

func Example_sumFloat32() {
	FromFloat32(1, 2, 3, 4.5).Sum().Println()
	// Output: 10.5
}
