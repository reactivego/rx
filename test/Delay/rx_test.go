package Delay

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_delay() {
	const ms = time.Millisecond
	start := time.Now()

	slice, err := FromInt(0, 1, 2, 3, 4).Delay(250 * ms).ToSlice()

	delay := time.Since(start)
	if 249*ms < delay && delay < 251*ms {
		fmt.Println("delay was 250 ms")
	}
	fmt.Println(slice)
	fmt.Println("error", err)
	// Output:
	// delay was 250 ms
	// [0 1 2 3 4]
	// error <nil>
}
