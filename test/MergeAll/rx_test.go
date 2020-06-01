package MergeAll

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_mergeAll() {
	source := CreateObservableString(func(N NextObservableString, E Error, C Complete, X Canceled) {
		for i := 0; i < 3; i++ {
			time.Sleep(100 * time.Millisecond)
			if X() {
				return
			}
			N(JustString(fmt.Sprintf("First %d", i)))
			N(JustString(fmt.Sprintf("Second %d", i)))
		}
		C()
	})

	source.MergeAll().Println()
	// Output:
	// First 0
	// Second 0
	// First 1
	// Second 1
	// First 2
	// Second 2
}

func Example_mergeAllInterval() {
	const ms = time.Millisecond

	intv1ms := Interval(1*ms).Take(2)
	intv10ms := Interval(10*ms).Take(2)
	intv100ms := Interval(100*ms).Take(2)

	FromObservableInt(intv1ms, intv10ms, intv100ms).MergeAll().Println()

	// Output:
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
}
