package Never

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_never() {
	concurrent := GoroutineScheduler()

	source := Never().Timeout(time.Hour).SubscribeOn(concurrent)

	never := true
	subscription := source.Subscribe(func(next interface{}, err error, done bool) {
		never = false
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		subscription.Unsubscribe()
	}()

	if subscription.Canceled() {
		fmt.Println("expected subscription to be active!")
	}
	subscription.Wait()
	if subscription.Subscribed() {
		fmt.Println("expected subscription to be canceled!")
	}
	if never && subscription.Canceled() {
		fmt.Println("OK")
	}
	// Output: OK
}
