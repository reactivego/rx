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

func ExampleObservable_Buffer() {
	const ms = time.Millisecond
	source := rx.Timer(0*ms, 100*ms).Take(4).ConcatMap(func(i interface{}) rx.Observable {
		switch i.(int) {
		case 0:
			return rx.From("a", "b")
		case 1:
			return rx.From("c", "d", "e")
		case 3:
			return rx.From("f", "g")
		}
		return rx.Empty()
	})
	closingNotifier := rx.Interval(100 * ms)
	source.Buffer(closingNotifier).Println()
	// Output:
	// [a b]
	// [c d e]
	// []
	// [f g]
}

func ExampleObservable_BufferTime() {
	const ms = time.Millisecond
	source := rx.Timer(0*ms, 100*ms).Take(4).ConcatMap(func(i interface{}) rx.Observable {
		switch i.(int) {
		case 0:
			return rx.From("a", "b")
		case 1:
			return rx.From("c", "d", "e")
		case 3:
			return rx.From("f", "g")
		}
		return rx.Empty()
	})
	source.BufferTime(100 * ms).Println()
	// Output:
	// [a b]
	// [c d e]
	// []
	// [f g]
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

func ExampleMerge() {
	a := rx.From(0, 2, 4)
	b := rx.From(1, 3, 5)
	rx.Merge(a, b).Println()
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleObservable_Distinct() {
	rx.From(1, 2, 2, 1, 3).Distinct().Println()

	// Output:
	// 1
	// 2
	// 3
}

func ExampleObservable_DistinctUntilChanged() {
	rx.From(1, 2, 2, 1, 3).DistinctUntilChanged().Println()

	// Output:
	// 1
	// 2
	// 1
	// 3
}

func ExampleObservable_MergeDelayError() {
	type any = interface{}
	const ms = time.Millisecond
	AddMul := func(add, mul int) func(any) any {
		return func(i any) any {
			return mul * (i.(int) + add)
		}
	}
	To := func(to int) func(any) any {
		return func(any) any {
			return to
		}
	}

	a := rx.Interval(20 * ms).Map(AddMul(1, 20)).Take(4).ConcatWith(rx.Throw(rx.RxError("boom")))
	b := rx.Timer(70*ms, 20*ms).Map(To(1)).Take(2)
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
	type any = interface{}
	const ms = time.Millisecond
	AddMul := func(add, mul int) func(any) any {
		return func(i any) any {
			return mul * (i.(int) + add)
		}
	}
	To := func(to int) func(any) any {
		return func(any) any {
			return to
		}
	}

	a := rx.Interval(20 * ms).Map(AddMul(1, 20)).Take(4).ConcatWith(rx.Throw(rx.RxError("boom")))
	b := rx.Timer(70*ms, 20*ms).Map(To(1)).Take(2)

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
			return rx.Range(n.(int), 2)
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

func ExampleObservable_MergeWith() {
	a := rx.From(0, 2, 4)
	b := rx.From(1, 3, 5)
	a.MergeWith(b).Println()
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
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

func ExampleObservable_StartWith() {
	rx.From(2, 3).StartWith(1).Println()

	// Output:
	// 1
	// 2
	// 3
}

func ExampleObservableObservable_SwitchAll() {
	type any = interface{}
	const ms = time.Millisecond

	// toObservable creates a new observable that emits an integer starting after 20ms and then repeated
	// every 20ms in the range starting at 0 and incrementing by 1. It takes only the first 10 emitted values.
	toObservable := func(i any) rx.Observable {
		return rx.Interval(20 * ms).Take(10)
	}

	rx.Interval(100 * ms).
		Take(3).
		MapObservable(toObservable).
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

func ExampleObservable_SubscribeOn_trampoline() {
	trampoline := rx.MakeTrampolineScheduler()
	observer := func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(trampoline.Count(), "print", next)
		case err != nil:
			fmt.Println(trampoline.Count(), "print", err)
		default:
			fmt.Println(trampoline.Count(), "print", "complete")
		}
	}
	fmt.Println(trampoline.Count(), "SUBSCRIBING...")
	subscription := rx.From(1, 2, 3).SubscribeOn(trampoline).Subscribe(observer)
	fmt.Println(trampoline.Count(), "WAITING...")
	subscription.Wait()
	fmt.Println(trampoline.Count(), "DONE")

	// Output:
	// 0 SUBSCRIBING...
	// 1 WAITING...
	// 1 print 1
	// 1 print 2
	// 1 print 3
	// 1 print complete
	// 0 DONE
}

func ExampleObservable_SubscribeOn_goroutine() {
	const ms = time.Millisecond
	goroutine := rx.GoroutineScheduler()
	observer := func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(goroutine.Count(), "print", next)
		case err != nil:
			fmt.Println(goroutine.Count(), "print", err)
		default:
			fmt.Println(goroutine.Count(), "print", "complete")
		}
	}
	fmt.Println(goroutine.Count(), "SUBSCRIBING...")
	subscription := rx.From(1, 2, 3).Delay(10 * ms).SubscribeOn(goroutine).Subscribe(observer)
	// Note that without a Delay the next Println lands at a random spot in the output.
	fmt.Println("WAITING...")
	subscription.Wait()
	fmt.Println(goroutine.Count(), "DONE")
	// Output:
	// 0 SUBSCRIBING...
	// WAITING...
	// 1 print 1
	// 1 print 2
	// 1 print 3
	// 1 print complete
	// 0 DONE
}

func ExampleObservable_WithLatestFrom() {
	a := rx.From(1, 2, 3, 4, 5)
	b := rx.From("A", "B", "C", "D", "E")
	a.WithLatestFrom(b).Println()
	// Output:
	// [2 A]
	// [3 B]
	// [4 C]
	// [5 D]
}

func ExampleObservableObservable_WithLatestFromAll() {
	a := rx.From(1, 2, 3, 4, 5)
	b := rx.From("A", "B", "C", "D", "E")
	c := rx.FromObservable(a, b)
	c.WithLatestFromAll().Println()
	// Output:
	// [2 A]
	// [3 B]
	// [4 C]
	// [5 D]
}
