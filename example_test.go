package x_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/x"
)

func Example_share() {
	serial := x.NewScheduler()

	shared := x.From(1, 2, 3).Share()

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
	serial := x.NewScheduler()

	// subject collects emits when there are no subscriptions active.
	in, out := x.Subject[int](0, 1)

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
		in.Error(x.Error("foo"))
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
	// foo
	// foo
}

func Example_multicast() {
	serial := x.NewScheduler()

	in, out := x.Multicast[int](1)

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
			in.Error(x.Error("foo"))
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
	// foo
	// foo
}

func Example_multicastDrop() {
	serial := x.NewScheduler()

	const onBackpressureDrop = -1

	// multicast with backpressure handling set to dropping incoming
	// items that don't fit in the buffer once it has filled up.
	in, out := x.Multicast[int](1 * onBackpressureDrop)

	// ignore everything before any subscriptions, including the last!
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Println().Go(serial)
	sub2 := out.Println().Go(serial)

	in.Next(2)               // accepted: buffer not full
	in.Next(3)               // dropped: buffer full
	in.Error(x.Error("foo")) // dropped: buffer full

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
	source := x.Empty[x.Observable[string]]()
	x.ConcatAll(source).Wait()

	source = x.Of(x.Empty[string]())
	x.ConcatAll(source).Wait()

	req := func(request string, duration time.Duration) x.Observable[string] {
		req := x.From(request + " response")
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

	source = x.From(req1).ConcatWith(x.From(req2, req3, req4).Delay(100 * ms))
	x.ConcatAll(source).Println().Wait()

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

	req := func(request string, duration time.Duration) x.Observable[string] {
		return x.From(request + " response").Delay(duration)
	}

	req1 := req("first", 50*ms)
	req2 := req("second", 10*ms)
	req3 := req("third", 60*ms)

	x.Race(req1, req2, req3).Println().Wait()

	err := func(text string, duration time.Duration) x.Observable[int] {
		return x.Throw[int](x.Error(text + " error")).Delay(duration)
	}

	err1 := err("first", 10*ms)
	err2 := err("second", 20*ms)
	err3 := err("third", 30*ms)

	fmt.Println(x.Race(err1, err2, err3).Wait(x.Goroutine))
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

	x.Map(x.Of(R{"Hello", "World"}).Marshal(json.Marshal), b2s).Println().Wait()
	// Output:
	// {"a":"Hello","b":"World"}
}

func Example_elementAt() {
	x.From(0, 1, 2, 3, 4).ElementAt(2).Println().Wait()
	// Output:
	// 2
}

func Example_exhaustAll() {
	const ms = time.Millisecond

	stream := func(name string, duration time.Duration, count int) x.Observable[string] {
		return x.Map(x.Timer[int](0*ms, duration), func(next int) string {
			return name + "-" + strconv.Itoa(next)
		}).Take(count)
	}

	streams := []x.Observable[string]{
		stream("a", 20*ms, 3),
		stream("b", 20*ms, 3),
		stream("c", 20*ms, 3),
		x.Empty[string](),
	}

	streamofstreams := x.Map(x.Timer[int](20*ms, 30*ms, 250*ms, 100*ms).Take(4), func(next int) x.Observable[string] {
		return streams[next]
	})

	err := x.ExhaustAll(streamofstreams).Println().Wait()

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
	source := x.From(0, 1, 2, 3)

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 1)")
	x.BufferCount(source, 2, 1).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 2)")
	x.BufferCount(source, 2, 2).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 3)")
	x.BufferCount(source, 2, 3).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 3, 2)")
	x.BufferCount(source, 3, 2).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 6, 6)")
	x.BufferCount(source, 6, 6).Println().Wait()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 0)")
	x.BufferCount(source, 2, 0).Println().Wait()
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

	interval42x4 := x.Interval[int](42 * ms).Take(4)
	interval16x4 := x.Interval[int](16 * ms).Take(4)

	err := x.SwitchAll(x.Map(interval42x4, func(next int) x.Observable[int] { return interval16x4 })).Println().Wait(x.Goroutine)

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

	webreq := func(request string, duration time.Duration) x.Observable[string] {
		return x.From(request + " result").Delay(duration)
	}

	first := webreq("first", 50*ms)
	second := webreq("second", 10*ms)
	latest := webreq("latest", 50*ms)

	switchmap := x.SwitchMap(x.Interval[int](20*ms).Take(3), func(i int) x.Observable[string] {
		switch i {
		case 0:
			return first
		case 1:
			return second
		case 2:
			return latest
		default:
			return x.Empty[string]()
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
	var first error = x.Error("error")
	a := x.Create(func(index int) (next int, err error, done bool) {
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
	source := x.From(1, 2, 3, 4, 5)

	count := source.Count()
	count.Println().Wait()

	emptySource := x.Empty[int]()
	emptyCount := emptySource.Count()
	emptyCount.Println().Wait()

	fmt.Println("OK")
	// Output:
	// 5
	// 0
	// OK
}

func Example_values() {
	source := x.From(1, 3, 5)

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
	source := x.From("ZERO", "ONE", "TWO")

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
	x.From(1, 2, 3, 4, 5).Skip(2).Println().Wait()
	// Output:
	// 3
	// 4
	// 5
}
