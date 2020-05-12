package CombineLatestMap

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	project := func(next int) ObservableString {
		return OfString(fmt.Sprintf(" %d ", next),fmt.Sprintf(">%d<", next))
	}

	FromInt(1, 2, 3, 4).CombineLatestMapString(project).Println()

	// Output:
	// [ 1   2   3   4 ]
	// [>1<  2   3   4 ]
	// [>1< >2<  3   4 ]
	// [>1< >2< >3<  4 ]
	// [>1< >2< >3< >4<]
}
