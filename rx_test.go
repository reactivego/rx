package rx_test

import (
	"fmt"
	"time"

	"github.com/reactivego/rx"
)

func ExampleObservable_All() {
	// Setup All to produce true only when all source values are less than 5
	lessthan5 := func(i interface{}) bool {
		return i.(int) < 5
	}

	result, err := rx.From(1, 2, 5, 2, 1).All(lessthan5).ToSingle()

	fmt.Println("All values less than 5?", result, err)

	result, err = rx.From(4, 1, 0, -1, 2, 3, 4).All(lessthan5).ToSingle()

	fmt.Println("All values less than 5?", result, err)
	// Output:
	// All values less than 5? false <nil>
	// All values less than 5? true <nil>
}

func ExampleObservable_ConcatWith() {
	oa := rx.From(0, 1, 2, 3)
	ob := rx.From(4, 5)
	oc := rx.From(6)
	od := rx.From(7, 8, 9)
	oa.ConcatWith(ob, oc).ConcatWith(od).Subscribe(func(next interface{}, err error, done bool) {
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

func ExampleDefer() {
	count := 0
	source := rx.Defer(func() rx.Observable {
		return rx.From(count)
	})
	mapped := source.Map(func(next interface{}) interface{} {
		return fmt.Sprintf("observable %d", next)
	})

	mapped.Println()
	count = 123
	mapped.Println()
	count = 456
	mapped.Println()

	// Output:
	// observable 0
	// observable 123
	// observable 456
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

func ExampleFrom_slice() {
	rx.From([]interface{}{1, 2, 3, 4, 5}...).Println()

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

func ExampleObservable_MergeDelayError() {
	const ms = time.Millisecond
	AddMul := func(add, mul int) func(interface{}) interface{}{
		return func(i interface{}) interface{} {
			return mul * (i.(int)+add)
		}
	}
	To := func(to int) func(interface{}) interface{} {
		return func(interface{}) interface{} {
			return to
		}
	}

	a := rx.Interval(20 * ms).AsObservable().Map(AddMul(1, 20)).Take(4).ConcatWith(rx.Throw(rx.RxError("boom")))
	b := rx.Timer(70 * ms, 20 * ms).AsObservable().Map(To(1)).Take(2)
	err := rx.MergeDelayError(a, b).Println()
	fmt.Println(err)
	// Output:
	// 20
	// 40
	// 60
	// 1
	// 80
	// 1
	// boom
}

func ExampleObservable_MergeDelayErrorWith() {
	const ms = time.Millisecond

	Access := func(slice ...interface{}) func(interface{}) interface{}{
		return func(i interface{}) interface{} {
			if i.(int) < len(slice) {
				return slice[i.(int)]
			} else {
				return rx.RxError("bool")
			}
		}
	}

	a := rx.Interval(20 * ms).AsObservable().Map(Access(20,40,60,80))
	b := rx.Timer(70 * ms, 20 * ms).AsObservable().Map(Access(1,1)).Take(2)

	fmt.Println(a.MergeDelayErrorWith(b).Println())
	// Output:
	// 20
	// 40
	// 60
	// 1
	// 80
	// 1
	// boom
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
