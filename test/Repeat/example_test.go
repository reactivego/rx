package Repeat

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	source := RepeatInt(5, 3)

	slice, err := source.ToSlice()

	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// [5 5 5]
	// <nil>
}

const _1M = 1000000

// Repeating the Observable 1 million times works fine using the
// Trampoline scheduler. When the repeated observable signals completion
// this will cause the Repeat operator to re-subscribe. The Trampoline
// scheduler will schedule every subscribe asynchronously to run after the
// first subscribe returns. So subscribe calls are actually not nested but
// executed in sequence.
func Example_trampoline() {
	source := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(1)
		N(2)
		N(3)
		C()
	})
	slice, err := source.Repeat(_1M).TakeLast(9).SubscribeOn(TrampolineScheduler()).ToSlice()

	fmt.Println(TrampolineScheduler())
	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// Trampoline{ Asynchronous:Serial(0) }
	// [1 2 3 1 2 3 1 2 3]
	// <nil>
}

// Repeating the Observable 1 million times works fine using the Goroutine
// scheduler. When the repeated observable signals completion, this will
// cause the Repeat operator to re-subscribe. The Goroutine scheduler
// will schedule every subscribe asynchronously and concurrently on a
// different goroutine. Therefore no nesting of subscriptions occurs.
// The Trampoline scheduler is much faster than the Goroutine
// scheduler because it does not create goroutines. It needs less than 70%
// of the time the Goroutine scheduler needs.
func Example_goroutine() {
	source := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(1)
		N(2)
		N(3)
		C()
	})
	slice, err := source.Repeat(_1M).TakeLast(9).SubscribeOn(GoroutineScheduler()).ToSlice()

	fmt.Println(GoroutineScheduler())
	fmt.Println(slice)
	fmt.Println(err)
	// Output:
	// Goroutine{ Asynchronous:Concurrent(0) }
	// [1 2 3 1 2 3 1 2 3]
	// <nil>
}
