package Timeout

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	start := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(1)
	source := CreateInt(func(observer IntObserver) {
		defer wg.Done()
		observer.Next(1)
		time.Sleep(time.Millisecond * 500)
		assert.True(t, observer.Closed(), "observer should be unsubscribed")
	})
	timed := source.Timeout(time.Millisecond * 250)

	actual, err := timed.ToSlice()
	expect := []int{1}
	assert.EqualError(t, err, ErrTimeout.Error())
	assert.True(t, err == ErrTimeout) // because const also identical to ErrTimeout

	elapsed := time.Now().Sub(start)
	assert.True(t, elapsed > time.Millisecond*250 && elapsed < time.Millisecond*500, "elapsed time should be between 250 and 500 ms")

	assert.Equal(t, expect, actual)

	wg.Wait()
}
