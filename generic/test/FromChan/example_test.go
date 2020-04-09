package FromChan

import "fmt"

func Example_fromChan() {
	ch := make(chan interface{}, 6)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	ch <- RxError("error")
	close(ch)

	err := FromChan(ch).AsObservableInt().Println()
	fmt.Println(err)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// error
}

func Example_fromChanInt() {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	
	err := FromChanInt(ch).Println()
	fmt.Println(err)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// <nil>
}
