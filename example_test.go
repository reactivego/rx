package rx_test

import (
	"fmt"
	"time"

	"github.com/reactivego/rx"
)

func ExampleObservable_Concat() {
	oa := rx.From(0, 1, 2, 3)
	ob := rx.From(4, 5)
	oc := rx.From(6)
	od := rx.From(7, 8, 9)
	oa.Concat(ob, oc).Concat(od).Subscribe(func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Printf("%d,", next.(int))
		case err != nil:
			fmt.Print("err", err)
		default:
			fmt.Printf("complete")
		}
	}).Wait()

	// Output:
	// 0,1,2,3,4,5,6,7,8,9,complete
}

func ExampleObservable_Do() {
	rx.From(1, 2, 3).Do(func(v interface{}) {
		fmt.Println(v.(int))
	}).Wait()

	// Output:
	// 1
	// 2
	// 3
}

func ExampleObservable_Filter() {
	even := func(i interface{}) bool {
		return i.(int)%2 == 0
	}

	rx.From(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).Println()

	// Output:
	// 2
	// 4
	// 6
	// 8
}

func ExampleFromChan() {
	ch := make(chan interface{}, 6)
	for i := 0; i < 5; i++ {
		ch <- i + 1
	}
	close(ch)

	rx.FromChan(ch).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleFrom() {
	rx.From(1, 2, 3, 4, 5).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleFromSlice() {
	rx.FromSlice([]interface{}{1, 2, 3, 4, 5}).Println()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleObservable_Map() {
	rx.From(1, 2, 3, 4).Map(func(i interface{}) interface{} {
		return fmt.Sprintf("%d!", i.(int))
	}).Println()

	// Output:
	// 1!
	// 2!
	// 3!
	// 4!
}

func ExampleObservable_MergeMap() {
	source := rx.From(1, 2).
		MergeMap(func(n interface{}) rx.Observable {
			return rx.Range(n.(int), 2).AsObservable()
		})
	if err := source.Println(); err != nil {
		panic(err)
	}

	// Output:
	// 1
	// 2
	// 2
	// 3
}

func ExampleObservable_Scan() {
	add := func(acc interface{}, value interface{}) interface{} {
		return acc.(int) + value.(int)
	}

	rx.From(1, 2, 3, 4, 5).Scan(add, 0).Println()

	// Output:
	// 1
	// 3
	// 6
	// 10
	// 15
}

func ExampleObservableObservable_SwitchAll() {
	// intToObs creates a new observable that emits an integer starting after and then repeated every 20 milliseconds
	// in the range starting at 0 and incrementing by 1. It takes only the first 10 emitted values and then uses
	// AsObservable to convert the IntObservable back to an untyped Observable.
	intToObs := func(i int) rx.Observable {
		return rx.Interval(20 * time.Millisecond).
			Take(10).
			AsObservable()
	}

	rx.Interval(100 * time.Millisecond).
		Take(3).
		MapObservable(intToObs).
		SwitchAll().
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
