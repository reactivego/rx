package PublishReplay

import (
	"fmt"
	"time"
)

func Example_channel() {
	ch := make(chan int, 1)
	ch <- 0
	go func() {
		for i := 1; i < 5; i++ {
			time.Sleep(50 * time.Millisecond)
			ch <- i
		}
		time.Sleep(100 * time.Millisecond)
		close(ch)
	}()
	source := FromChanInt(ch).PublishReplay(0, 500*time.Millisecond)
	source.Connect()

	if err := source.Wait(); err != nil {
		fmt.Println(err)
	}

	print := func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("complete")
			}
		}
	}

	sub := source.Subscribe(print)
	if sub.Closed() {
		fmt.Println("subscription was closed")
	} else {
		fmt.Println("subscription was still subscribed")
	}

	time.Sleep(300 * time.Millisecond)

	sub = source.Subscribe(print)
	if sub.Closed() {
		fmt.Println("subscription was closed")
	} else {
		fmt.Println("subscription was still subscribed")
	}

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// complete
	// subscription was closed
	// 3
	// 4
	// complete
	// subscription was closed
}
