package Average

func Example_averageInt() {
	FromInts(1, 2, 3, 4, 5).Average().Println()

	// Output:
	// 3
}

func Example_averageFloat32() {
	FromFloat32s(1, 2, 3, 4).Average().Println()

	// Output:
	// 2.5
}
