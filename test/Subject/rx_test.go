package Subject

import (
	"fmt"
	"sort"
	"sync"

	_ "github.com/reactivego/rx/generic"
)

// Subject is both an observer and an observable.
func Example_subject() {
	subject := NewSubjectInt()
	scheduler := GoroutineScheduler()

	subscription := subject.SubscribeOn(scheduler).Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		}
	})

	subject.Next(123)
	subject.Next(456)
	subject.Complete()

	subscription.Wait()
	// Output:
	// 123
	// 456
}

// Subject will forward errors from the observer to the observable side.
func Example_subjectError() {
	subject := NewSubjectInt()

	// feed the subject...
	go func() {
		subject.Error(RxError("something bad happened"))
	}()

	err := subject.Wait()
	fmt.Println(err)
	// Output:
	// something bad happened
}

// Subject has an observable side that provides multicasting. This means that
// two subscribers will receive the same data at approximately the same time.
func Example_subjectMultiple() {
	subject := NewSubjectInt()
	scheduler := GoroutineScheduler()

	var messages struct {
		sync.Mutex
		values []string
	}

	var subscriptions []Subscription
	for i := 0; i < 5; i++ {
		index := i
		subscription := subject.SubscribeOn(scheduler).Subscribe(func(next int, err error, done bool) {
			if !done {
				message := fmt.Sprint(index, next)

				messages.Lock()
				messages.values = append(messages.values, message)
				messages.Unlock()
			}
		})
		subscriptions = append(subscriptions, subscription)
	}

	subject.Next(123)
	subject.Next(456)
	subject.Complete()
	subject.Next(111)
	subject.Next(222)

	for i := 0; i < 5; i++ {
		subscriptions[i].Wait()
	}

	sort.Sort(sort.StringSlice(messages.values))
	for _, message := range messages.values {
		fmt.Println(message)
	}

	// Output:
	// 0 123
	// 0 456
	// 1 123
	// 1 456
	// 2 123
	// 2 456
	// 3 123
	// 3 456
	// 4 123
	// 4 456
}
