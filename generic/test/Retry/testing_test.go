package Retry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryCurrentGoroutine(t *testing.T) {
	errored := false
	a := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			// Error triggers subscribe and subscribe is scheduled on trampoline....
			observer.Error(RxError("error"))
			errored = true
		}
	}).SubscribeOn(CurrentGoroutineScheduler())
	b, e := a.Retry().ToSlice()
	assert.NoError(t, e)
	assert.Equal(t, []int{1, 2, 3, 1, 2, 3}, b)
	assert.True(t, errored)
}

func TestRetryNewGoroutine(t *testing.T) {
	errored := false
	a := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			observer.Error(RxError("error"))
			errored = true
		}
	}).SubscribeOn(NewGoroutineScheduler())
	b, e := a.Retry().ToSlice()
	assert.NoError(t, e)
	assert.Equal(t, []int{1, 2, 3, 1, 2, 3}, b)
	assert.True(t, errored)
}
