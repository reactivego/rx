package SwitchMap

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_switchMap() {
	delay := func(duration time.Duration, value Observable) Observable {
		return Never().Timeout(duration).Catch(value)
	}

	webreq := func(request string, duration time.Duration) ObservableString {
		return delay(duration, From(request+" result")).AsObservableString()
	}

	first := webreq("first", 50*time.Millisecond)
	second := webreq("second", 10*time.Millisecond)
	latest := webreq("latest", 50*time.Millisecond)

	err := Interval(20 * time.Millisecond).Take(3).SwitchMapString(func(i int) ObservableString {
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
