package Publish

import (
	"fmt"
	"runtime"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	// source is an ObservableInt
	source := FromInt(1, 2)

	// pub is a ConnectableInt
	pub := source.Publish()

	// scheduler will run a task asynchronously on a new goroutine.
	scheduler := GoroutineScheduler()

	// First subscriber (asynchronous)
	sub1 := pub.SubscribeOn(scheduler).
		MapString(func(v int) string {
			return fmt.Sprintf("value is %d", v)
		}).
		Subscribe(func(next string, err error, done bool) {
			if !done {
				fmt.Println("sub1", next)
			}
		})

	// Second subscriber (asynchronous)
	sub2 := pub.SubscribeOn(scheduler).
		MapBool(func(v int) bool {
			return v > 1
		}).
		Subscribe(func(next bool, err error, done bool) {
			if !done {
				fmt.Println("sub2", next)
			}
		})

	// Third subscriber (asynchronous)
	sub3 := pub.SubscribeOn(scheduler).
		Subscribe(func(next int, err error, done bool) {
			if !done {
				fmt.Println("sub3", next)
			}
		})

	// Connect will cause the publisher to subscribe to the source
	pub.Connect()

	// Wait for all subscribers to finish.
	sub1.Wait()
	sub2.Wait()
	sub3.Wait()

	// Unordered output:
	// sub1 value is 1
	// sub1 value is 2
	// sub2 false
	// sub2 true
	// sub3 1
	// sub3 2
}

func Example_publishRefCount() {
	scheduler := GoroutineScheduler()
	ch := make(chan int, 30)
	s := FromChanInt(ch).Publish().RefCount().SubscribeOn(scheduler)
	a := []int{}
	b := []int{}
	appendToSlice := func(slice *[]int) IntObserveFunc {
		return func(next int, err error, done bool) {
			if !done {
				*slice = append(*slice, next)
			}
		}
	}
	asub := s.Subscribe(appendToSlice(&a))
	bsub := s.Subscribe(appendToSlice(&b))
	ch <- 1
	ch <- 2
	ch <- 3
	// make sure the channel gets enough time to be fully processed.
	for i := 0; i < 10; i++ {
		time.Sleep(20 * time.Millisecond)
		runtime.Gosched()
	}
	asub.Unsubscribe()
	if asub.Subscribed() {
		fmt.Println("asub should be closed")
	}
	ch <- 4
	close(ch)
	bsub.Wait()

	fmt.Println("a", a)
	fmt.Println("b", b)

	// Output:
	// a [1 2 3]
	// b [1 2 3 4]
}
