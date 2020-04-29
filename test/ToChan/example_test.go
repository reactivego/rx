package ToChan

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_basic() {
	ch := From(1, 2, "middle", 2, 1).ToChan()
	for v := range ch {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// middle
	// 2
	// 1
}

func Example_unsubscribeDirect() {
	sub := NewSubscriber()
	ch := From(1, 2, "middle", 2, 1).ToChan(sub)
	if i, ok := <-ch; ok {
		fmt.Println(i)
	}
	sub.Unsubscribe()
	// if i, ok := <-ch; ok {
	// 	fmt.Println(i)
	// }
	// if i, ok := <-ch; ok {
	// 	fmt.Println(i)
	// }
	fmt.Println("unsubscribed:", sub.Canceled())

	// Output:
	// 1
	// unsubscribed: true
}

func Example_unsubscribeLoop() {
	sub := NewSubscriber()
	ch := From(1, 2, "middle", 2, 1).ToChan(sub)
	for i := range ch {
		switch v := i.(type) {
		case int:
			fmt.Println(v)
		case string:
			sub.Unsubscribe()
		}
	}
	fmt.Println("unsubscribed:", sub.Canceled())

	// Output:
	// 1
	// 2
	// unsubscribed: true
}

// The source ObservableInt created by Range(1,9) will emit values in the range
// 1 to 9. ToChan subscribes asynchronously to the source which begins to emit
// the values into the channel. The channel itself is retured by ToChan and can
// be used by the caller in e.g. a for loop. ToChan changes the scheduler it
// uses for subscribing from the default synchronous trampoline to an
// asynchronous one. The subscription runs concurrently with the channel reading
// for loop.
func Example_intChannel() {
	for value := range Range(1, 9).ToChan() {
		fmt.Print(value)
	}

	// Output: 123456789
}

// First a source Observable is created that emits 1 to 8 followed by an error.
// ToChan will take the error and output it as the last item emitted by the
// channel.
func Example_inlineError() {
	source := Create(func(observer Observer) {
		for i := 1; i < 9; i++ {
			observer.Next(i)
		}
		observer.Error(RxError("sad"))
	})

	for value := range source.ToChan() {
		switch value.(type) {
		case int:
			fmt.Print(value)
		case error:
			fmt.Printf("<%v>", value)
		}
	}

	// Output: 12345678<sad>
}
