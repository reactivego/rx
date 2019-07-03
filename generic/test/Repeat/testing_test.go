package Repeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	source := RepeatInt(5, 3)
	expect := []int{5, 5, 5}
	result, err := source.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
	result, err = source.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

const _1M = 1000000

// Repeating the Observable 1 million times works fine using the
// CurrentGoroutine scheduler. When the repeated observable signals completion
// this will cause the Repeat operator to re-subscribe. The CurrentGoroutine
// scheduler will schedule every subscribe asynchronously to run after the
// first subscribe returns. So subscribe calls are actually not nested but
// executed in sequence.
func TestRepeatCurrentGoroutine(t *testing.T) {
	source := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		observer.Complete()
	})
	result, err := source.Repeat(_1M).TakeLast(9).ToSlice(SubscribeOn(CurrentGoroutineScheduler()))
	expect := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

// Repeating the Observable 1 million times works fine using the NewGoroutine
// scheduler. When the repeated observable signals completion, this will
// cause the Repeat operator to re-subscribe. The NewGoroutine scheduler
// will schedule every subscribe asynchronously and concurrently on a
// different goroutine. Therefore no nesting of subscriptions occurs.
// The CurrentGoroutine scheduler is much faster than the NewGoroutine
// scheduler because it does not create goroutines. It needs less than 70% 
// of the time the NewGoroutine scheduler needs.
func TestRepeatNewGoroutineScheduler(t *testing.T) {
	source := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		observer.Complete()
	})
	result, err := source.Repeat(_1M).TakeLast(9).ToSlice(SubscribeOn(NewGoroutineScheduler()))
	expect := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
