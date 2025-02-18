package SwitchMap

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_switchMap() {
	const ms = time.Millisecond

	webreq := func(request string, duration time.Duration) ObservableString {
		return FromString(request + " result").Delay(duration)
	}

	first := webreq("first", 50*ms)
	second := webreq("second", 10*ms)
	latest := webreq("latest", 50*ms)

	err := IntervalInt(20 * ms).Take(3).SwitchMapString(func(i int) ObservableString {
		switch i {
		case 0:
			return first
		case 1:
			return second
		case 2:
			return latest
		default:
			return EmptyString()
		}
	}).Println()
	fmt.Println(err)
	// Output:
	// second result
	// latest result
	// <nil>
}
