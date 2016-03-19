// +build !special

package rxgo

import (
	"fmt"
	"rxgo/observable"
	"time"
)

type Unsubber interface {
	Unsub()
	Unsubbed() bool
}

type Unsub struct {
	unsub func()
}

func (f *Unsub) Unsub() {
	if f.unsub != nil {
		f.unsub()
		f.unsub = nil
	}
}
func (f *Unsub) Unsubbed() bool {
	return f.unsub == nil
}

func MakeUnsubber() Unsubber {
	return &Unsub{func() { fmt.Println("test") }}
}

func ExampleUnsubber() {
	u := MakeUnsubber()
	u.Unsub()
	u.Unsub()
	u.Unsub()
	u.Unsub()
	fmt.Println(u.Unsubbed())

	//Output:
	// test
	// true
}

func ExampleRange() {
	done := make(chan struct{})
	observable.Range(1, 10).Subscribe(func(next int, err error, completed bool) {
		if err != nil || completed {
			close(done)
			fmt.Println("completed")
		} else {
			fmt.Println(next)
		}
	})
	<-done

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
	// completed
}

func ExampleReplay() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(time.Millisecond * 100)
		}
		close(ch)
	}()
	done := make(chan struct{})
	s := observable.FromIntChannel(ch).Replay(0, time.Millisecond*600).Subscribe(func(next int, err error, completed bool) {
		switch {
		case err != nil:
			println(err)
			close(done)
		case completed:
			fmt.Println("completed")
			close(done)
		default:
			fmt.Println(next)
		}
	})

	<-done

	time.Sleep(time.Millisecond * 500)
	if s.Unsubscribed() {
		fmt.Println("s was unsubscribed")
	} else {
		fmt.Println("s was still subscribed")
	}

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// completed
	// s was unsubscribed
}

func ExampleMapString() {
	a := observable.FromInts(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToArray()
	for _, v := range a {
		fmt.Println(v)
	}

	//Output:
	// 1!
	// 2!
	// 3!
	// 4!
}

func ExampleSubscribeOn() {
	done := make(chan struct{})
	unsubscriber := observable.FromInts(1, 2, 3, 4, 5).SubscribeOn(observable.GoroutineScheduler).Subscribe(func(next int, err error, completed bool) {
		switch {
		case err != nil:
			fmt.Println(err)
		case completed:
			fmt.Println("completed")
			close(done)
		default:
			fmt.Println(next)
		}
	})

	// Check that Executor does not leak out via Unsubscriber returned by Subscribe function.
	if executor, ok := unsubscriber.(observable.Executor); ok {
		fmt.Println("Executor interface was exported by Subscribe")
		fmt.Println(executor)
	} else {
		fmt.Println("Executor interface was NOT exported by Subscribe")
	}
	<-done

	//Output:
	// Executor interface was NOT exported by Subscribe
	// 1
	// 2
	// 3
	// 4
	// 5
	// completed
}
