package MapTo

import (
	_ "github.com/reactivego/rx"
)

func Example_mapToString() {
	FromInt(1, 2, 3, 4).MapToString("YES").Println()

	//Output:
	// YES
	// YES
	// YES
	// YES
}
