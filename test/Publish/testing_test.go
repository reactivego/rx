package Publish

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPublishRefCount(t *testing.T) {
	scheduler := NewGoroutine()
	ch := make(chan int, 30)
	s := FromChanInt(ch).Publish().RefCount().SubscribeOn(scheduler)
	a := []int{}
	b := []int{}
	asub := s.SubscribeNext(func(n int) { a = append(a, n) }, SubscribeOn(scheduler))
	bsub := s.SubscribeNext(func(n int) { b = append(b, n) }, SubscribeOn(scheduler))
	ch <- 1
	ch <- 2
	ch <- 3
	// make sure the channel gets enough time to be fully processed.
	for i := 0; i < 10; i++ {
		time.Sleep(20 * time.Millisecond)
		runtime.Gosched()
	}
	asub.Unsubscribe()
	assert.True(t, asub.Closed())
	ch <- 4
	close(ch)
	bsub.Wait()
	assert.Equal(t, []int{1, 2, 3}, a)
	assert.Equal(t, []int{1, 2, 3, 4}, b)
	//assert.True(t, false, "force fail")
}
