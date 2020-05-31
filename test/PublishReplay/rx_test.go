package PublishReplay

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

// Basic example with the maximum buffer capacity of approx. 32000 items where
// items are retained forever.
func Example_basic() {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	source := FromChanInt(ch).PublishReplay(0, 0)
	source.Connect().Wait()

	err := source.Println()
	fmt.Println(err)

	err = source.Println()
	fmt.Println(err)
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// <nil>
	// 0
	// 1
	// 2
	// 3
	// 4
	// <nil>
}

// Example that shows using a buffer that retains only the latest 2 values and
// the use of AutoConnect to declaratively call the Connect method.
func Example_replayWithSize() {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	concurrent := GoroutineScheduler()
	source := FromChanInt(ch).PublishReplay(2, 0).AutoConnect(1).SubscribeOn(concurrent)

	err := source.Println()
	fmt.Println(err)

	err = source.Println()
	fmt.Println(err)

	err = source.Println()
	fmt.Println(err)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// <nil>
	// 3
	// 4
	// <nil>
	// 3
	// 4
	// <nil>
}

// Example that shows how items may be retained for only a limited time.
func Example_replayWithExpiry() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()
	source := FromChanInt(ch).PublishReplay(0, 600*time.Millisecond)

	scheduler := GoroutineScheduler()
	scheduler.Schedule(func() {

		time.Sleep(500 * time.Millisecond)
		// 500ms has passed, everything should still be present
		err := source.Println()
		fmt.Println(err)

		time.Sleep(100 * time.Millisecond)

		// 600ms has passed, first value should be gone by now.
		err = source.Println()
		fmt.Println(err)

	})

	source.Connect().Wait()

	scheduler.Wait()

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// <nil>
	// 1
	// 2
	// 3
	// 4
	// <nil>
}

func Example_replayWithExpiryAndSubscribe() {
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

	concurrent := GoroutineScheduler()
	concurrent.Schedule(func(){
		if err := source.Wait(); err != nil {
			fmt.Println(err)
		}
	})

	source.Connect().Wait()
	concurrent.Wait()

	print := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	sub := source.Subscribe(print)
	sub.Wait()
	if sub.Subscribed() {
		fmt.Println("subscription was still subscribed")
	} else {
		fmt.Println("subscription was unsubscribed")
	}

	time.Sleep(300 * time.Millisecond)

	sub = source.Subscribe(print)
	sub.Wait()
	if sub.Subscribed() {
		fmt.Println("subscription was still subscribed")
	} else {
		fmt.Println("subscription was unsubscribed")
	}

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// complete
	// subscription was unsubscribed
	// 3
	// 4
	// complete
	// subscription was unsubscribed
}
