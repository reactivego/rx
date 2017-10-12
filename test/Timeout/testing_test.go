package Timeout

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1)
	actual, err := CreateInt(func(observer IntObserver) {
		defer wg.Done()
		observer.Next(1)
		time.Sleep(time.Millisecond * 500)
		assert.True(t, observer.Closed(), "observer should be unsubscribed")
	}).
		Timeout(time.Millisecond * 250).
		ToSlice()
	elapsed := time.Now().Sub(start)
	assert.Error(t, err)
	assert.Equal(t, ErrTimeout, err)
	assert.True(t, elapsed > time.Millisecond*250 && elapsed < time.Millisecond*500, "elapsed time should be between 250 and 500 ms")
	assert.Equal(t, []int{1}, actual)
	wg.Wait()
}
