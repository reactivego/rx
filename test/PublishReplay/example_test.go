package PublishReplay

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_channel() {
	ch := make(chan int, 1)
	ch <- 0
	go func() {
		for i := 1; i < 5; i++ {
			time.Sleep(50 * time.Millisecond)
			ch <- i
		}
		time.Sleep(100 * time.Millisecond)
		close(ch)
	}()
	source := FromChanInt(ch).PublishReplay(0, 500*time.Millisecond)
	source.Connect()

	if err := source.Wait(); err != nil {
		fmt.Println(err)
	}

	print := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	sub := source.Subscribe(print)
	sub.Wait()
	if sub.Canceled() {
		fmt.Println("subscription was canceled")
	} else {
		fmt.Println("subscription was still subscribed")
	}

	time.Sleep(300 * time.Millisecond)

	sub = source.Subscribe(print)
	sub.Wait()
	if sub.Canceled() {
		fmt.Println("subscription was canceled")
	} else {
		fmt.Println("subscription was still subscribed")
	}

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// complete
	// subscription was canceled
	// 3
	// 4
	// complete
	// subscription was canceled
}

/*
// Basic example with the maximum buffer capacity of approx. 32000 items where
// items are retained forever.
func TestReplayBasic(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	source := FromChanInt(ch).PublishReplay(0, 0)
	source.Connect()

	expected := []int{0, 1, 2, 3, 4}

	result, err := source.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	result, err = source.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// Example that shows using a buffer that retains only the latest 2 values and
// the use of AutoConnect to declaratively call the Connect method.
func TestReplayWithSize(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	source := FromChanInt(ch).PublishReplay(2, 0).AutoConnect(1)
	result, err := source.ToSlice()
	expect := []int{0, 1, 2, 3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	result, err = source.ToSlice()
	expect = []int{3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	result, err = source.ToSlice()
	expect = []int{3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

// Example that shows how items may be retained for only a limited time.
func TestReplayWithExpiry(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()
	source := FromChanInt(ch).PublishReplay(0, 600*time.Millisecond)
	source.Connect()

	time.Sleep(500 * time.Millisecond)

	// 500ms has passed, everything should still be present
	result, err := source.ToSlice()
	expect := []int{0, 1, 2, 3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	time.Sleep(100 * time.Millisecond)

	// 600ms has passed, first value should be gone by now.
	result, err = source.ToSlice()
	expect = []int{1, 2, 3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
*/
