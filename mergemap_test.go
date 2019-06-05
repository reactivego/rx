package rx

func ExampleObservable_MergeMap() {
	scheduler := NewGoroutineScheduler()

	if err := From(1, 2).MergeMap(func(n any) Observable {
		return Range(n.(int), 2).AsObservable()
	}).Println(SubscribeOn(scheduler)); err != nil {
		panic(err)
	}

	// Unordered output:
	// 1
	// 2
	// 2
	// 3
}
