package Distinct

func Example_distinct() {
	FromInts(1, 1, 2, 2, 3, 2, 4, 5).Distinct().Println()
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
