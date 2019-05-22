package Debounce

import (
	"fmt"
	"time"
)

func Example_debounce() {

	source := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		time.Sleep(300 * time.Millisecond)

		observer.Next(2)
		time.Sleep(80 * time.Millisecond) // 80ms < 100ms => '2' is ignored

		observer.Next(3)
		time.Sleep(110 * time.Millisecond)

		observer.Next(4)
		observer.Complete()
	})

	debounced := source.Debounce(100 * time.Millisecond)

	debounced.SubscribeNext(func(next int) { fmt.Println(next) })

	// Output:
	// 1
	// 3
	// 4
}
