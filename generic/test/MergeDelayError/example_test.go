package MergeDelayError

import (
	"fmt"
	"time"
)

func Example_mergeDelayError() {
	const _5ms = 5 * time.Millisecond
	const _10ms = 10 * time.Millisecond

	sourceA := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Error(RxError("error.sourceA"))
	})

	sourceB := CreateInt(func(observer IntObserver) {
		time.Sleep(_5ms)
		observer.Next(0)
		time.Sleep(_10ms)
		observer.Next(2)
		observer.Complete()
	})

	result, err := sourceA.MergeDelayError(sourceB).ToSlice()
	fmt.Println(result, err)

	// Output: [1 0 2] error.sourceA
}
