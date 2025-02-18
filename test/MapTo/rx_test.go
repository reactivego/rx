package MapTo

import (
	_ "github.com/reactivego/rx/generic"
)

func Example_mapToString() {
	FromInt(1, 2, 3, 4).MapToString("YES").Println()

	//Output:
	// YES
	// YES
	// YES
	// YES
}
