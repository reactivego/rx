package Delay

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_delay() {
	const (
		_240ms = 240 * time.Millisecond
		_250ms = 250 * time.Millisecond
		_260ms = 260 * time.Millisecond
	)

	start := time.Now()

	slice, err := FromInt(1, 2, 3).Delay(_250ms).ToSlice()

	// check the slice was created after a delay of 250 milliseconds
	duration := time.Since(start)
	if duration < _240ms || duration > _260ms {
		fmt.Println("delay must be between 240 and 260 miliseconds")
	} else {
		fmt.Println("delay was 250 miliseconds")
	}

	// print the slice we got
	fmt.Println(slice)

	// print any error that was returned
	fmt.Println("error", err)
	// Output:
	// delay was 250 miliseconds
	// [1 2 3]
	// error <nil>
}
