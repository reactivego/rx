package SwitchAll

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSwitchAllGoroutine(t *testing.T) {
	scheduler := NewGoroutineScheduler()
	items,err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(10)
		}).
		SwitchAll().
		ToSlice(SubscribeOn(scheduler))

	expect := []int{0,1, 0,1, 0,1, 0,1,2,3,4,5,6,7,8,9}
	assert.NoError(t, err)
	assert.Equal(t, expect, items)
}

func TestSwitchAllTrampoline(t *testing.T) {
	scheduler := CurrentGoroutineScheduler()
	items, err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(10)
		}).
		SwitchAll().
		ToSlice(SubscribeOn(scheduler))

	expect := []int{0,1,2,3,4,5,6,7,8,9}
	assert.NoError(t, err)
	assert.Equal(t, expect, items)
}
