// go test -run=XXX -bench=rx -cpu=1,2,3,4,5,6,7,8 -timeout=1h -count=10
package ObserverObservable

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reactivego/scheduler"

	_ "github.com/reactivego/rx/generic"
)

func Example() {
	concurrent := scheduler.Goroutine

	observer, observable := MakeObserverObservable(0, 128)

	print := func(next interface{}, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	observable.SubscribeOn(concurrent).Subscribe(print)
	observable.SubscribeOn(concurrent).Subscribe(print)
	observable.SubscribeOn(concurrent).Subscribe(print)

	observer.Next(1)
	for start := time.Now(); time.Since(start) < 1*time.Millisecond; {
		runtime.Gosched()
	}
	observer.Next(2)
	observer.Complete()

	concurrent.Wait()
	// Output:
	// 1
	// 1
	// 1
	// 2
	// 2
	// 2
	// complete
	// complete
	// complete
}

func TestLazySender(t *testing.T) {
	concurrent := scheduler.Goroutine

	const ms = time.Millisecond
	const SUBSCRIBERS = 2

	observer, observable := MakeObserverObservable(0, 10, 10, SUBSCRIBERS)

	concurrent.ScheduleFuture(245*ms, func() {
		for i := 0; i < 100; i++ {
			observer.Next(i)
		}
		observer.Complete()
	})

	for i := 0; i < SUBSCRIBERS; i++ {
		concurrent.Schedule(func() {
			subscription := observable.Subscribe(func(next interface{}, err error, done bool) {
				switch {
				case !done:
					//
				case err != nil:
					//
				default:
					//
				}
			})
			subscription.Wait()
		})
	}

	concurrent.Wait()
}

func TestSleepingReceiver(t *testing.T) {
	concurrent := scheduler.Goroutine

	observer, observable := MakeObserverObservable(0, 2, 2, 1)

	subscriber := observable.SubscribeOn(concurrent).Subscribe(func(next interface{}, err error, done bool) {

	})

	time.Sleep(300 * time.Millisecond)
	observer.Next(1)
	observer.Complete()

	err := subscriber.Wait()

	if err != nil {
		t.Error(err)
	}

	concurrent.Wait()
}

func TestChanMaxAge(t *testing.T) {
	const ms = time.Millisecond

	observer, observable := MakeObserverObservable(50*ms, 128)

	// emit entry every 1ms; 0 at 0ms, 1 at 1ms and so on to 99 at 99ms
	start := time.Now()
	for i := 0 * ms; i < 100; i++ {
		// wait i milliseconds
		for time.Since(start) < i*ms {
			runtime.Gosched()
		}
		observer.Next(int(i))
	}
	observer.Complete()

	num := 50
	subscription := observable.Subscribe(func(next interface{}, err error, done bool) {
		if !done {
			if next.(int) != num {
				t.Fatalf("expected %d, got %d", num, next)
			}
			num++
		}
	})
	subscription.Wait()
}

func TestChanEndpointKeep(t *testing.T) {
	age := 0 * time.Millisecond
	len, cap, ecap := 0, 128, 32
	observer, observable := MakeObserverObservable(age, len, cap, ecap)

	for i := 0; i < 100; i++ {
		observer.Next(i)
	}
	observer.Complete()

	num := 0
	subscription := observable.Subscribe(func(next interface{}, err error, done bool) {
		if !done {
			num++
		}
	})
	subscription.Wait()
	if num != 0 {
		t.Fatal("Got", num, "buffered values but I ask for none (keep arg was 0)")
	}
}

func Test_1xN(t *testing.T) {
	concurrent := scheduler.Goroutine

	const SUBSCRIBERS = 16

	observer, observable := MakeObserverObservable(0, 10, 10, SUBSCRIBERS)

	var subscribing sync.WaitGroup
	subscribing.Add(SUBSCRIBERS)
	concurrent.Schedule(func() {
		subscribing.Wait()
		for i := 0; i < 100; i++ {
			observer.Next(i)
		}
		observer.Complete()
	})

	for i := 0; i < SUBSCRIBERS; i++ {
		concurrent.Schedule(func() {
			expected := 4950
			total := 0
			count := 0
			last := -1
			subscription := observable.Subscribe(func(next interface{}, err error, done bool) {
				switch {
				case !done:
					current := next.(int)
					if last != current-1 {
						fmt.Fprintln(os.Stderr, "last", last, "current", current)
					}
					last = current
					total += last
					count++
				case err != nil:
					fmt.Fprintln(os.Stderr, err)
				default:
					if total != expected {
						fmt.Fprintln(os.Stderr, "complete", "last", last, "count", count, "total", total)
					}
				}
			})
			subscribing.Done()
			subscription.Wait()
			if total != expected {
				t.Error(fmt.Sprint(i), ":", total, " != ", expected)
			}
		})
	}

	concurrent.Wait()
}

func Benchmark_1xN_rx(b *testing.B) {
	const BUFSIZE = 500
	CPU := runtime.GOMAXPROCS(0)
	// fmt.Fprintln(os.Stderr, "NUM =", b.N, "; CPU =", CPU)

	observer, observable := MakeObserverObservable(0, BUFSIZE)

	var subscribing sync.WaitGroup
	subscribing.Add(CPU)

	var receiving sync.WaitGroup
	receiving.Add(CPU)

	var sending sync.WaitGroup
	sending.Add(1)

	go func() {
		subscribing.Wait()

		start := time.Now()

		// send b.N+1 items to work around internal benchmark graining. Graining
		// prevents all pb.Next() calls completing because pb.Next() tokens are
		// still present in observables that completed.
		for i := 0; i < b.N+1; i++ {
			observer.Next(i)
		}
		observer.Complete()

		nps := time.Now().Sub(start).Nanoseconds() / int64(b.N+1)
		// b.Logf("1x%d, %d msg(s), %d ns/send, %.1fM msgs/sec", CPU, b.N+1, nps, 1.0e03/float64(nps))
		_ = nps

		receiving.Wait()

		sending.Done()
	}()

	scount, rcount := int64(0), int64(0)

	channelindex := int64(-1)
	b.RunParallel(func(pb *testing.PB) {
		pbnext := func() bool {
			atomic.AddInt64(&scount, 1)
			if pb.Next() {
				atomic.AddInt64(&rcount, 1)
				return true
			}
			// stop timing the other messages arriving.
			b.StopTimer()
			return false
		}
		ci := atomic.AddInt64(&channelindex, 1)
		var sum, count int64
		subscription := observable.Subscribe(func(next interface{}, err error, done bool) {
			switch {
			case !done:
				if pbnext() {
					var expectedSum int64
					sum += int64(next.(int))
					count++
					expectedSum = int64(count) * int64(count-1) / 2
					if sum != expectedSum {
						b.Fatalf("data corruption at count == %d ; expected == %d got sum == %d", count, expectedSum, sum)
					}
				}
			case err != nil:
				fmt.Fprintln(os.Stderr, ci, err)
			default:
				// fmt.Fprintln(os.Stderr, ci, "COMPLETE")
			}
		})
		subscribing.Done()
		subscription.Wait()
		receiving.Done()
	})
	// time the remainder of the sender completing
	// b.StartTimer()

	sending.Wait()

	if atomic.LoadInt64(&rcount) != int64(b.N) {
		b.Errorf("data loss; expected %d messages got %d (%d)", b.N, atomic.LoadInt64(&rcount), atomic.LoadInt64(&scount))
	}
}

func Benchmark_1xN_go(b *testing.B) {
	const BUFSIZE = 500
	CPU := runtime.GOMAXPROCS(0)
	// fmt.Fprintln(os.Stderr, "N =", b.N, "; CPU =", CPU)

	// create channels
	var channels []chan interface{}
	for p := 0; p < CPU; p++ {
		channels = append(channels, make(chan interface{}, BUFSIZE))
	}

	var subscribing sync.WaitGroup
	subscribing.Add(CPU)

	var receiving sync.WaitGroup
	receiving.Add(CPU)

	var sending sync.WaitGroup
	sending.Add(1)

	go func() {
		// wait for all receivers to subscribe
		subscribing.Wait()

		start := time.Now()

		for i := 0; i < b.N+1; i++ {
			// To make comparison fair we also have go channel calculate time on send.
			// The mcast channel also does this once for every Send call.
			updated := time.Since(start).Nanoseconds()
			_ = updated
			for _, c := range channels {
				c <- i
			}
		}

		// close channels
		for _, c := range channels {
			close(c)
		}

		nps := time.Now().Sub(start).Nanoseconds() / int64(b.N+1)
		// b.Logf("1x%d, %d msg(s), %d ns/send, %.1fM msgs/sec", CPU, b.N+1, nps, 1.0e03/float64(nps))
		_ = nps

		// wait for receivers
		receiving.Wait()

		sending.Done()
	}()

	scount, rcount := int64(0), int64(0)
	channelindex := int64(-1)
	b.RunParallel(func(pb *testing.PB) {
		pbnext := func() bool {
			atomic.AddInt64(&scount, 1)
			if pb.Next() {
				atomic.AddInt64(&rcount, 1)
				return true
			}
			// stop timing the other messages arriving.
			b.StopTimer()
			return false
		}
		ci := atomic.AddInt64(&channelindex, 1)
		ch := channels[ci]
		subscribing.Done()
		var sum, count int64
		for next := range ch {
			if pbnext() {
				sum += int64(next.(int))
				count++
				expectedSum := count * (count - 1) / 2
				if sum != expectedSum {
					b.Errorf("data corruption; at count %d, expected sum %d got %d", count, expectedSum, sum)
				}
			}
		}
		receiving.Done()
	})
	// time the remainder of the sender completing
	// b.StartTimer()

	// You only get here after the benchmark is done!!!!!
	sending.Wait()

	if rcount != int64(b.N) {
		b.Errorf("data loss; expected %d messages got %d (%d)", b.N, atomic.LoadInt64(&rcount), atomic.LoadInt64(&scount))
	}
}
