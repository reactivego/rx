package Buffer

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_buffer() {
	const ms = time.Millisecond

	source := TimerInt(0*ms, 100*ms).Take(4).ConcatMap(func(i int) Observable {
		switch i {
		case 0:
			return From("a", "b")
		case 1:
			return From("c", "d", "e")
		case 2:
			return Empty()
		case 3:
			return From("f", "g")
		default:
			return Empty()
		}
	}).Publish().AutoConnect(1)

	source.Buffer(Interval(100 * ms)).Println()

	// Output:
	// [a b]
	// [c d e]
	// []
	// [f g]
}
