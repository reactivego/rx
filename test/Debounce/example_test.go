package Debounce

import (
	"fmt"
	"time"
)

func Example_debounce() {

	source := CreateInt(func(observer IntObserver) {
		time.Sleep(100 * time.Millisecond)
		observer.Next(1)
		time.Sleep(300 * time.Millisecond)
		observer.Next(2)
		time.Sleep(80 * time.Millisecond) // duration too short, Next(2) will be ignored.
		observer.Next(3)
		time.Sleep(110 * time.Millisecond)
		observer.Next(4)
		observer.Complete()
	})

	source.Debounce(100 * time.Millisecond).Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		}
	})

	// Output:
	// 1
	// 3
	// 4
}
