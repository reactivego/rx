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

// Repeating the Observable 1 million times works fine using the Trampoline
// scheduler. When the repeated observable signals completion this will
// cause the Repeat operator to re-subscribe. The Trampoline scheduler
// will schedule every subscribe asynchronously to run after the first
// subscribe returns. So subscribe calls are actually not nested but
// executed in sequence.
func TestRepeatTrampoline(t *testing.T) {
	source := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		observer.Complete()
	})
	result, err := source.Repeat(_1M).TakeLast(9).ToSlice(SubscribeOn(NewTrampolineScheduler()))
	expect := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

// Repeating the Observable 1 million times works fine using the Goroutine
// scheduler. When the repeated observable signals completion this will
// cause the Repeat operator to re-subscribe. The Goroutine scheduler
// will schedule every subscribe asynchronously and concurrently on a
// different goroutine. Therefore no nesting of subscriptions occurs.
// The Trampoline (because it does not create goroutines) needs less
// than 70% of the time the Goroutine scheduler needs.
func TestRepeatGoroutine(t *testing.T) {
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
