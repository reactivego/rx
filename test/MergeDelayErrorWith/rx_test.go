package MergeDelayErrorWith

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_mergeDelayErrorWith() {
	const _5ms = 5 * time.Millisecond
	const _10ms = 10 * time.Millisecond

	sourceA := CreateInt(func(N NextInt, E Error, _ Complete, _ Canceled) {
		N(1)
		E(RxError("error.sourceA"))
	})

	sourceB := CreateInt(func(N NextInt, _ Error, C Complete, _ Canceled) {
		time.Sleep(_5ms)
		N(0)
		time.Sleep(_10ms)
		N(2)
		C()
	})

	result, err := sourceA.MergeDelayErrorWith(sourceB).ToSlice()
	fmt.Println(result, err)
	// Output: [1 0 2] error.sourceA
}
