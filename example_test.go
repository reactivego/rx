package rx

import (
	"fmt"
	"time"

)

func ExampleObservable_Concat() {
	oa := From(0, 1, 2, 3)
	ob := From(4, 5)
	oc := From(6)
	od := From(7, 8, 9)
	oa.Concat(ob, oc).Concat(od).SubscribeNext(func(next any) {
		fmt.Printf("%d,", next.(int))
	})

	// Output:
	// 0,1,2,3,4,5,6,7,8,9,
}

func ExampleObservable_Do() {
	From(1, 2, 3).Do(func(v any) {
		fmt.Println(v.(int))
	}).Wait()

	// Output:
	// 1
	// 2
	// 3
}

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

func ExampleObservable_Map() {
	From(1, 2, 3, 4).Map(func(i any) any {
		return fmt.Sprintf("%d!", i.(int))
	}).Println()

	// Output:
	// 1!
	// 2!
	// 3!
	// 4!
}

func ExampleObservable_MergeMap() {
	scheduler := NewGoroutineScheduler()

	source := From(1, 2).
		MergeMap(func(n any) Observable {
			return Range(n.(int), 2).AsObservable()
		}).
		SubscribeOn(scheduler)
	if err := source.Println(); err != nil {
		panic(err)
	}

	// Unordered output:
	// 1
	// 2
	// 2
	// 3
}

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

func ExampleObservableObservable_SwitchAll() {

	// SwitchAll does not work well with the default trampoline scheduler, so we use a goroutine scheduler instead.
	scheduler := NewGoroutineScheduler()

	// intToObs creates a new observable that emits an integer starting after and then repeated every 20 milliseconds
	// in the range starting at 0 and incrementing by 1. It takes only the first 10 emitted values and then uses
	// AsObservable to convert the IntObservable back to an untyped Observable.
	intToObs := func(i int) Observable {
		return Interval(20 * time.Millisecond).
			Take(10).
			AsObservable()
	}

	Interval(100 * time.Millisecond).
		Take(3).
		MapObservable(intToObs).
		SwitchAll().
		SubscribeOn(scheduler).
		Println()

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 0
	// 1
	// 2
	// 3
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}
