package TakeWhile

func Example_takeWhile() {
	FromInt(1, 2, 3, 4, 5).TakeWhile(func (next int) bool { return next < 4 }).Println()

	// Output:
	// 1
	// 2
	// 3
}
