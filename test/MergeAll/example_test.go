package MergeAll

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_mergeAll() {
	source := CreateObservableString(func(observer ObservableStringObserver) {
		for i := 0; i < 3; i++ {
			time.Sleep(100 * time.Millisecond)
			observer.Next(JustString(fmt.Sprintf("First %d", i)))
			observer.Next(JustString(fmt.Sprintf("Second %d", i)))
		}
		observer.Complete()
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
