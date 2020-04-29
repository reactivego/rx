package Debounce

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_debounce() {

	i := 1
	due := []time.Duration{
		0,
		300 * time.Millisecond,
		80 * time.Millisecond, // 80ms < 100ms => '2' is ignored
		110 * time.Millisecond,
		0,
	}
	source := MakeTimedInt(due[0],
		func(Next func(int), Error func(error), Complete func()) time.Duration {
			if i < len(due) {
				Next(i)
				i++
				return due[i-1]
			} else {
				Complete()
				return 0
			}
		})

	debounced := source.Debounce(100 * time.Millisecond)

	if err := debounced.Println(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1
	// 3
	// 4
}
