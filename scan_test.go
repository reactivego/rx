package rx

func ExampleObservable_Scan() {
	add := func(acc any, value any) any {
		return acc.(int) + value.(int)
	}

	From(1, 2, 3, 4, 5).Scan(add, 0).Println()

	// Output:
	// 1
	// 3
	// 6
	// 10
	// 15
}
