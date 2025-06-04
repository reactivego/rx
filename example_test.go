package rx_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/reactivego/rx"
	"github.com/reactivego/scheduler"
)

func Example_share() {
	serial := rx.NewScheduler()

	shared := rx.From(1, 2, 3).Share()

	shared.Println().Go(serial)
	shared.Println().Go(serial)
	shared.Println().Go(serial)

	serial.Wait()
	// Output:
	// 1
	// 1
	// 1
	// 2
	// 2
	// 2
	// 3
	// 3
	// 3
}

func Example_subject() {
	serial := rx.NewScheduler()

	// subject collects emits when there are no subscriptions active.
	in, out := rx.Subject[int](0, 1)

	// ignore everything before any subscriptions, except the last because buffer size is 1
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Println().Go(serial)
	sub2 := out.Println().Go(serial)

	// schedule the subsequent emits on the serial scheduler otherwise these calls
	// will block because the buffer is full.
	// subject will detect usage of scheduler on observable side and use it on the
	// observer side to keep the data flow through the subject going.
	serial.Schedule(func() {
		in.Next(2)
		in.Next(3)
		in.Done(rx.Err)
	})

	serial.Wait()
	fmt.Println(sub1.Wait())
	fmt.Println(sub2.Wait())
	// Output:
	// 1
	// 1
	// 2
	// 2
	// 3
	// 3
	// rx
	// rx
}

func Example_multicast() {
	serial := rx.NewScheduler()

	in, out := rx.Multicast[int](1)

	// Ignore everything before any subscriptions, including the last!
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// Schedule the subsequent emits in a loop. This will be the first task to
	// run on the serial scheduler after the subscriptions have been added.
	serial.ScheduleLoop(2, func(index int, again func(next int)) {
		if index < 4 {
			in.Next(index)
			again(index + 1)
		} else {
			in.Done(rx.Err)
		}
	})

	// Add a couple of subscriptions
	sub1 := out.Println().Go(serial)
	sub2 := out.Println().Go(serial)

	// Let the scheduler run and wait for all of its scheduled tasks to finish.
	serial.Wait()
	fmt.Println(sub1.Wait())
	fmt.Println(sub2.Wait())
	// Output:
	// 2
	// 2
	// 3
	// 3
	// rx
	// rx
}

func Example_multicastDrop() {
	serial := rx.NewScheduler()

	const onBackpressureDrop = -1

	// multicast with backpressure handling set to dropping incoming
	// items that don't fit in the buffer once it has filled up.
	in, out := rx.Multicast[int](1 * onBackpressureDrop)

	// ignore everything before any subscriptions, including the last!
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Println().Go(serial)
	sub2 := out.Println().Go(serial)

	in.Next(2)      // accepted: buffer not full
	in.Next(3)      // dropped: buffer full
	in.Done(rx.Err) // dropped: buffer full

	serial.Wait()
	fmt.Println(sub1.Wait())
	fmt.Println(sub2.Wait())
	// Output:
	// 2
	// 2
	// <nil>
	// <nil>
}

func Example_concatAll() {
	source := rx.Empty[rx.Observable[string]]()
	rx.ConcatAll(source).Wait()

	source = rx.Of(rx.Empty[string]())
	rx.ConcatAll(source).Wait()

	req := func(request string, duration time.Duration) rx.Observable[string] {
		req := rx.From(request + " response")
		if duration == 0 {
			return req
		}
		return req.Delay(duration)
	}

	const ms = time.Millisecond

	req1 := req("first", 10*ms)
	req2 := req("second", 20*ms)
	req3 := req("third", 0*ms)
	req4 := req("fourth", 60*ms)

	source = rx.From(req1).ConcatWith(rx.From(req2, req3, req4).Delay(100 * ms))
	rx.ConcatAll(source).Println().Wait()

	fmt.Println("OK")
	// Output:
	// first response
	// second response
	// third response
	// fourth response
	// OK
}

func Example_race() {
	const ms = time.Millisecond

	req := func(request string, duration time.Duration) rx.Observable[string] {
		return rx.From(request + " response").Delay(duration)
	}

	req1 := req("first", 50*ms)
	req2 := req("second", 10*ms)
	req3 := req("third", 60*ms)

	rx.Race(req1, req2, req3).Println().Wait()

	err := func(text string, duration time.Duration) rx.Observable[int] {
		return rx.Throw[int](errors.New(text + " error")).Delay(duration)
	}

	err1 := err("first", 10*ms)
	err2 := err("second", 20*ms)
	err3 := err("third", 30*ms)

	fmt.Println(rx.Race(err1, err2, err3).Wait(rx.Goroutine))
	// Output:
	// second response
	// first error
}

func Example_marshal() {
	type R struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	b2s := func(data []byte) string { return string(data) }

	rx.Map(rx.Of(R{"Hello", "World"}).Marshal(json.Marshal), b2s).Println().Wait()
	// Output:
	// {"a":"Hello","b":"World"}
}

func Example_elementAt() {
	rx.From(0, 1, 2, 3, 4).ElementAt(2).Println().Wait()
	// Output:
	// 2
}

func Example_exhaustAll() {
	const ms = time.Millisecond

	stream := func(name string, duration time.Duration, count int) rx.Observable[string] {
		return rx.Map(rx.Timer[int](0*ms, duration), func(next int) string {
			return name + "-" + strconv.Itoa(next)
		}).Take(count)
	}

	streams := []rx.Observable[string]{
		stream("a", 20*ms, 3),
		stream("b", 20*ms, 3),
		stream("c", 20*ms, 3),
		rx.Empty[string](),
	}

	streamofstreams := rx.Map(rx.Timer[int](20*ms, 30*ms, 250*ms, 100*ms).Take(4), func(next int) rx.Observable[string] {
		return streams[next]
	})

	err := rx.ExhaustAll(streamofstreams).Println().Wait()

	if err == nil {
		fmt.Println("success")
	}
	// Output:
	// a-0
	// a-1
	// a-2
	// c-0
	// c-1
	// c-2
	// success
}

func Example_bufferCount() {
	source := rx.From(0, 1, 2, 3)

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 1)")
	rx.BufferCount(source, 2, 1).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 2)")
	rx.BufferCount(source, 2, 2).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 3)")
	rx.BufferCount(source, 2, 3).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 3, 2)")
	rx.BufferCount(source, 3, 2).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 6, 6)")
	rx.BufferCount(source, 6, 6).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 0)")
	rx.BufferCount(source, 2, 0).Println().Wait()
	// Output:
	// BufferCount(From(0, 1, 2, 3), 2, 1)
	// [0 1]
	// [1 2]
	// [2 3]
	// [3]
	// BufferCount(From(0, 1, 2, 3), 2, 2)
	// [0 1]
	// [2 3]
	// BufferCount(From(0, 1, 2, 3), 2, 3)
	// [0 1]
	// [3]
	// BufferCount(From(0, 1, 2, 3), 3, 2)
	// [0 1 2]
	// [2 3]
	// BufferCount(From(0, 1, 2, 3), 6, 6)
	// [0 1 2 3]
	// BufferCount(From(0, 1, 2, 3), 2, 0)
	// [0 1]
}

func Example_switchAll() {
	const ms = time.Millisecond

	interval42x4 := rx.Interval[int](42 * ms).Take(4)
	interval16x4 := rx.Interval[int](16 * ms).Take(4)

	err := rx.SwitchAll(rx.Map(interval42x4, func(next int) rx.Observable[int] { return interval16x4 })).Println().Wait(rx.Goroutine)

	if err == nil {
		fmt.Println("success")
	}
	// Output:
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 2
	// 3
	// success
}

func Example_switchMap() {
	const ms = time.Millisecond

	webreq := func(request string, duration time.Duration) rx.Observable[string] {
		return rx.From(request + " result").Delay(duration)
	}

	first := webreq("first", 50*ms)
	second := webreq("second", 10*ms)
	latest := webreq("latest", 50*ms)

	switchmap := rx.SwitchMap(rx.Interval[int](20*ms).Take(3), func(i int) rx.Observable[string] {
		switch i {
		case 0:
			return first
		case 1:
			return second
		case 2:
			return latest
		default:
			return rx.Empty[string]()
		}
	})

	err := switchmap.Println().Wait()
	if err == nil {
		fmt.Println("success")
	}
	// Output:
	// second result
	// latest result
	// success
}

func Example_retry() {
	var first error = rx.Err
	a := rx.Create(func(index int) (next int, err error, done bool) {
		if index < 3 {
			return index, nil, false
		}
		err, first = first, nil
		return 0, err, true
	})
	err := a.Retry().Println().Wait()
	fmt.Println(first == nil)
	fmt.Println(err)
	// Output:
	// 0
	// 1
	// 2
	// 0
	// 1
	// 2
	// true
	// <nil>
}

func Example_count() {
	source := rx.From(1, 2, 3, 4, 5)

	count := source.Count()
	count.Println().Wait()

	emptySource := rx.Empty[int]()
	emptyCount := emptySource.Count()
	emptyCount.Println().Wait()

	fmt.Println("OK")
	// Output:
	// 5
	// 0
	// OK
}

func Example_values() {
	source := rx.From(1, 3, 5)

	// Why choose the Goroutine concurrent scheduler?
	// An observable can actually be at the root of a tree
	// of separately running observables that have their
	// responses merged. The Goroutine scheduler allows
	// these observables to run concurrently.

	// run the observable on 1 or more goroutines
	for i := range source.Values(scheduler.Goroutine) {
		// This is called from a newly created goroutine
		fmt.Println(i)
	}

	// run the observable on the current goroutine
	for i := range source.Values(scheduler.New()) {
		fmt.Println(i)
	}

	fmt.Println("OK")
	// Output:
	// 1
	// 3
	// 5
	// 1
	// 3
	// 5
	// OK
}

func Example_all() {
	source := rx.From("ZERO", "ONE", "TWO")

	for k, v := range source.All() {
		fmt.Println(k, v)
	}

	fmt.Println("OK")
	// Output:
	// 0 ZERO
	// 1 ONE
	// 2 TWO
	// OK
}

func Example_skip() {
	rx.From(1, 2, 3, 4, 5).Skip(2).Println().Wait()
	// Output:
	// 3
	// 4
	// 5
}

func Example_autoConnect() {
	// Create a multicaster hot observable that will emit every 100 milliseconds
	hot := rx.Interval[int](100 * time.Millisecond).Take(10).Publish()
	hotsub := hot.Connect(rx.Goroutine)
	defer hotsub.Unsubscribe()
	fmt.Println("Hot observable created and emitting 0,1,2,3,4,5,6 ...")

	// Publish the hot observable again but only Connect to it when 2
	// subscribers have connected.
	source := hot.Take(5).Publish().AutoConnect(2)

	// First subscriber
	sub1 := source.Printf("Subscriber 1: %d\n").Go()
	fmt.Println("First subscriber connected, waiting a bit...")

	// Wait a bit, nothing will emit yet
	time.Sleep(525 * time.Millisecond)

	fmt.Println("Second subscriber connecting, emissions begin!")
	// Second subscriber triggers the connection
	sub2 := source.Printf("Subscriber 2: %d\n").Go()

	// Wait for emissions to complete
	hotsub.Wait()
	sub1.Wait()
	sub2.Wait()

	// Unordered output:
	// Hot observable created and emitting 0,1,2,3,4,5,6 ...
	// First subscriber connected, waiting a bit...
	// Second subscriber connecting, emissions begin!
	// Subscriber 1: 5
	// Subscriber 2: 5
	// Subscriber 1: 6
	// Subscriber 2: 6
	// Subscriber 1: 7
	// Subscriber 2: 7
	// Subscriber 1: 8
	// Subscriber 2: 8
	// Subscriber 1: 9
	// Subscriber 2: 9
}

func Example_mergeMap() {
	source := rx.From("https://reactivego.io", "https://github.com/reactivego")

	merged := rx.MergeMap(source, func(next string) rx.Observable[string] {
		fakeFetchData := rx.Of(fmt.Sprintf("content of %q", next))
		return fakeFetchData
	})

	merged.Println().Go().Wait()

	// Output:
	// content of "https://reactivego.io"
	// content of "https://github.com/reactivego"
}
