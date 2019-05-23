package rx

func ExampleFromChan() {
	ch := make(chan any, 6)
	for i := 0; i < 5; i++ {
		ch <- i + 1
	}
	close(ch)

	FromChan(ch).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleFrom() {
	From(1, 2, 3, 4, 5).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleFromSlice() {
	FromSlice([]any{1, 2, 3, 4, 5}).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
