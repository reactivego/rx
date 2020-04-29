package Count

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_count() {
	count := FromInt(1, 2, 3, 4, 5, 6, 7).Count()

	err := count.Println()
	if err != nil {
		fmt.Println("error =", err)
	}

	// resubscribe, expect that it will return 7 again
	result, err := count.ToSingle()
	fmt.Println(result)
	if err != nil {
		fmt.Println("error =", err)
	}

	// Output:
	// 7
	// 7
}
