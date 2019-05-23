package rx

func ExampleObservable_Filter() {
	even := func(i any) bool {
	 return i.(int)%2 == 0
	}

	From(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).Println()

	// Output:
	// 2
	// 4
	// 6
	// 8
}
