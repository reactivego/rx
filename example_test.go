package observable_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "github.com/reactivego/observable"
)

func Example_subject() {
	serial := NewScheduler()

	// subject collects emits when there are no subscriptions active.
	in, out := Subject[int](0, 1)

	// ignore everything before any subscriptions, except the last because buffer size is 1
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Subscribe(Println[int](), serial)
	sub2 := out.Subscribe(Println[int](), serial)

	// schedule the subsequent emits on the serial scheduler otherwise these calls
	// will block because the buffer is full.
	// subject will detect usage of scheduler on observable side and use it on the
	// observer side to keep the data flow through the subject going.
	serial.Schedule(func() {
		in.Next(2)
		in.Next(3)
		in.Error(Error("foo"))
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
	serial := NewScheduler()

	in, out := Multicast[int](1)

	// ignore everything before any subscriptions, including the last!
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Subscribe(Println[int](), serial)
	sub2 := out.Subscribe(Println[int](), serial)

	// schedule the subsequent emits on a separate goroutine, otherwise they'll block
	go func() {
		in.Next(2)
		in.Next(3)
		in.Error(Error("foo"))
	}()

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
	serial := NewScheduler()

	const onBackpressureDrop = -1

	// multicast with backpressure handling set to dropping incoming
	// items that don't fit in the buffer once it has filled up.
	in, out := Multicast[int](1 * onBackpressureDrop)

	// ignore everything before any subscriptions, including the last!
	in.Next(-2)
	in.Next(-1)
	in.Next(0)
	in.Next(1)

	// add a couple of subscriptions
	sub1 := out.Subscribe(Println[int](), serial)
	sub2 := out.Subscribe(Println[int](), serial)

	in.Next(2)             // accepted: buffer not full
	in.Next(3)             // dropped: buffer full
	in.Error(Error("foo")) // dropped: buffer full

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
	source := Empty[Observable[string]]()
	ConcatAll(source).Wait()

	source = Of(Empty[string]())
	ConcatAll(source).Wait()

	req := func(request string, duration time.Duration) Observable[string] {
		req := From(request + " response")
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

	source = From(req1).ConcatWith(From(req2, req3, req4).Delay(100 * ms))
	ConcatAll(source).Println()

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

	req := func(request string, duration time.Duration) Observable[string] {
		return From(request + " response").Delay(duration)
	}

	req1 := req("first", 50*ms)
	req2 := req("second", 10*ms)
	req3 := req("third", 60*ms)

	Race(req1, req2, req3).Println()

	err := func(text string, duration time.Duration) Observable[int] {
		return Throw[int](Error(text + " error")).Delay(duration)
	}

	err1 := err("first", 10*ms)
	err2 := err("second", 20*ms)
	err3 := err("third", 30*ms)

	fmt.Println(Race(err1, err2, err3).Wait(Goroutine))
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

	Map(Of(R{"Hello", "World"}).Marshal(json.Marshal), b2s).Println()
	// Output:
	// {"a":"Hello","b":"World"}
}

func Example_elementAt() {
	From(0, 1, 2, 3, 4).ElementAt(2).Println()
	// Output:
	// 2
}

func Example_exhaustAll() {
	const ms = time.Millisecond

	stream := func(name string, duration time.Duration, count int) Observable[string] {
		return Map(Timer[int](0*ms, duration), func(next int) string {
			return name + "-" + strconv.Itoa(next)
		}).Take(count)
	}

	streams := []Observable[string]{
		stream("a", 20*ms, 3),
		stream("b", 20*ms, 3),
		stream("c", 20*ms, 3),
		Empty[string](),
	}

	streamofstreams := Map(Timer[int](20*ms, 30*ms, 250*ms, 100*ms).Take(4), func(next int) Observable[string] {
		return streams[next]
	})

	err := ExhaustAll(streamofstreams).Println()

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
	source := From(0, 1, 2, 3)

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 1)")
	BufferCount(source, 2, 1).Println()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 2)")
	BufferCount(source, 2, 2).Println()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 3)")
	BufferCount(source, 2, 3).Println()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 3, 2)")
	BufferCount(source, 3, 2).Println()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 6, 6)")
	BufferCount(source, 6, 6).Println()

	fmt.Println("BufferCount(From(0, 1, 2, 3), 2, 0)")
	BufferCount(source, 2, 0).Println()

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

	interval42x4 := Interval[int](42 * ms).Take(4)
	interval16x4 := Interval[int](16 * ms).Take(4)

	err := SwitchAll(Map(interval42x4, func(next int) Observable[int] { return interval16x4 })).Println(Goroutine)

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

	webreq := func(request string, duration time.Duration) Observable[string] {
		return From(request + " result").Delay(duration)
	}

	first := webreq("first", 50*ms)
	second := webreq("second", 10*ms)
	latest := webreq("latest", 50*ms)

	switchmap := SwitchMap(Interval[int](20*ms).Take(3), func(i int) Observable[string] {
		switch i {
		case 0:
			return first
		case 1:
			return second
		case 2:
			return latest
		default:
			return Empty[string]()
		}
	})

	err := switchmap.Println()
	if err == nil {
		fmt.Println("success")
	}
	// Output:
	// second result
	// latest result
	// success
}

func Example_retry() {
	var first error = Error("error")
	a := Create(func(index int) (next int, err error, done bool) {
		if index < 3 {
			return index, nil, false
		}
		err, first = first, nil
		return 0, err, true
	})
	err := a.Retry().Println()
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
