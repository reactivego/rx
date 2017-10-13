package Publish

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFork(t *testing.T) {
	ch := make(chan int, 30)
	s := FromChanInt(ch).Share() // allready does a subscribe, but nothing in channel yet...
	a := []int{}
	b := []int{}
	asub := s.SubscribeNext(func(n int) { a = append(a, n) })
	bsub := s.SubscribeNext(func(n int) { b = append(b, n) })
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
	assert.Equal(t, []int{1, 2, 3, 4}, b)
	assert.Equal(t, []int{1, 2, 3}, a)
	//assert.True(t, false, "force fail")
}
