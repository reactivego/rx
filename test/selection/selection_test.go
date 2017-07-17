package selection

import (
	"errors"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestSample(t *testing.T) {
	a := Interval(time.Millisecond * 90).Sample(time.Millisecond * 200).Take(3).ToArray()
	assert.Equal(t, []int{1, 3, 5}, a)
}

func TestDebounce(t *testing.T) {
	s := CreateInt(func(observer IntSubscriber) {
		time.Sleep(100 * time.Millisecond)
		observer.Next(1)
		time.Sleep(300 * time.Millisecond)
		observer.Next(2)
		time.Sleep(80 * time.Millisecond)
		observer.Next(3)
		time.Sleep(110 * time.Millisecond)
		observer.Next(4)
		observer.Complete()
	})
	a := s.Debounce(time.Millisecond * 100).ToArray()
	assert.Equal(t, []int{1, 3, 4}, a)
}

func TestConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	oa := FromIntArray(a)
	ob := FromIntArray(b)
	oc := FromIntArray(c)
	s := oa.Concat(ob).Concat(oc).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	s = oa.Concat(ob, oc).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
}

func TestMerge(t *testing.T) {
	sa := CreateInt(func(observer IntSubscriber) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		time.Sleep(10 * time.Millisecond)
		observer.Next(3)
		observer.Complete()
	})
	sb := CreateInt(func(observer IntSubscriber) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.Merge(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2, 3}, a)
}

func TestMergeDelayError(t *testing.T) {
	sa := CreateInt(func(observer IntSubscriber) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		observer.Error(errors.New("error"))
	})
	sb := CreateInt(func(observer IntSubscriber) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.MergeDelayError(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2}, a)
}

func TestRecover(t *testing.T) {
	// Create three different observables.
	// They are all passive until someone subscribes to them...
	o123 := FromInts(1, 2, 3)
	o45 := FromInts(4, 5)
	oThrowError := ThrowInt(errors.New("error"))

	a, err := o123.Concat(oThrowError).Catch(o45).ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestRetry(t *testing.T) {
	errored := false
	a := CreateInt(func(observer IntSubscriber) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			observer.Error(errors.New("error"))
			errored = true
		}
	})
	b := a.Retry().ToArray()
	assert.Equal(t, []int{1, 2, 3, 1, 2, 3}, b)
	assert.True(t, errored)
}

func TestTimeout(t *testing.T) {
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1)
	actual, err := CreateInt(func(subscriber IntSubscriber) {
		subscriber.Next(1)
		time.Sleep(time.Millisecond * 500)
		assert.True(t, subscriber.Unsubscribed())
		wg.Done()
	}).
		Timeout(time.Millisecond * 250).
		ToArrayWithError()
	elapsed := time.Now().Sub(start)
	assert.Error(t, err)
	assert.Equal(t, ErrTimeout, err)
	assert.True(t, elapsed > time.Millisecond*250 && elapsed < time.Millisecond*500)
	assert.Equal(t, []int{1}, actual)
	wg.Wait()
}

func TestFork(t *testing.T) {
	ch := make(chan int, 30)
	s := FromIntChannel(ch).Fork() // allready does a subscribe, but channel blocks.
	a := []int{}
	b := []int{}
	sub := s.SubscribeNext(func(n int) { a = append(a, n) })
	s.SubscribeNext(func(n int) { b = append(b, n) })
	ch <- 1
	ch <- 2
	ch <- 3
	// make sure the channel gets enough time to be fully processed.
	for i := 0; i < 10; i++ {
		time.Sleep(20 * time.Millisecond)
		runtime.Gosched()
	}
	sub.Unsubscribe()
	assert.True(t, sub.Unsubscribed())
	ch <- 4
	close(ch)
	s.Wait()
	assert.Equal(t, []int{1, 2, 3, 4}, b)
	assert.Equal(t, []int{1, 2, 3}, a)
}

func TestFlatMap(t *testing.T) {
	actual, err := Range(1, 2).FlatMap(func(n int) ObservableInt { return Range(n, 2) }).ToArrayWithError()
	assert.NoError(t, err)
	sort.Ints(actual)
	assert.Equal(t, []int{1, 2, 2, 3}, actual)
}
