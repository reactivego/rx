package All

import _ "github.com/reactivego/rx"

func Example_all() {
	LessThan5 := func(i int) bool {
		return i < 5
	}
	FromInt(1, 2, 6, 2, 1).All(LessThan5).Println("All values less than 5?")
	// Output: All values less than 5? false
}
