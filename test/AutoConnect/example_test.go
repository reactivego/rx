package AutoConnect

import "fmt"

// This example shows how to use AutoConnect with a PublishReplay connectable.
// The first Subscribe call will cause Connect on PublishReplay so it subscribes
// to range1to9. The second Subscribe call will cause the sequence to be
// replayed without doing a subscribe to range1to9.
func Example_autoConnect() {

	range1to9 := DeferInt(func () ObservableInt {
		fmt.Println("subscribed")
		i := 1
		return CreateRecursiveInt(func(next NextInt, error Error, complete Complete) {
			if i < 10 {
				next(i)
				i++
			} else {
				complete()
			}
		})
	})

	source := range1to9.PublishReplay(10, 0).AutoConnect(1)

	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Print(next, ",")
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	fmt.Println("first...")
	source.Subscribe(observe).Wait()
	fmt.Println("second...")
	source.Subscribe(observe).Wait()

	// Output:
	// first...
	// subscribed
	// 1,2,3,4,5,6,7,8,9,complete
	// second...
	// 1,2,3,4,5,6,7,8,9,complete
}

// This example how to use an asynchronous scheduler, so multiple subscribe
// calls can be made to an AutoConnect(2) operator that will only connect on
// the second subscription.
func Example_autoConnectMulti() {
	// We will be subscribing on an asynchronous scheduler, otherwise first
	// call to Subscribe will lockup.
	scheduler := GoroutineScheduler()

	range1to99 := CreateInt(func(next NextInt, error Error, complete Complete, canceled Canceled) {
		fmt.Println("subscribed")
		for i := 1; i < 100; i++ {
			if canceled() {
				return
			}
			next(i)
		}
		complete()
	})

	source := range1to99.PublishReplay(10, 0).AutoConnect(2).SubscribeOn(scheduler)

	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			// ignore values.
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	fmt.Println("first...")

	// Subcribe is asynchronous now
	sub1 := source.Subscribe(observe)

	fmt.Println("second...")

	// Also asynchronous
	sub2 := source.Subscribe(observe)

	fmt.Println("wait...")

	// We now need to wait for the subscriptions to terminate.
	sub1.Wait()
	sub2.Wait()

	fmt.Println("done")

	// Output:
	// first...
	// second...
	// wait...
	// subscribed
	// complete
	// complete
	// done
}
