package DoOnComplete

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_doOnComplete() {
	complete := false
	result, err := EmptyInt().DoOnComplete(func() { complete = true }).ToSlice()
	fmt.Println(result, err, complete)
	// Output: [] <nil> true
}

func Example_doOnCompleteSubscribe() {
	wait := make(chan struct{})
	source := FromInt(1, 2, 3, 4, 5).DoOnComplete(func() { close(wait) })

	result := []int{}
	concurrent := GoroutineScheduler()
	source.SubscribeOn(concurrent).Subscribe(func(next int, err error, done bool) {
		if !done {
			result = append(result, next)
		}
	})

	<-wait
	fmt.Println(result)
	// Output: [1 2 3 4 5]
}
