package Debounce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {
	source := CreateInt(func(observer IntObserver) {
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
	result, err := source.Debounce(time.Millisecond * 100).ToSlice()
	expect := []int{1, 3, 4}

	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
