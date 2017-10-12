package Merge

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	sourceA := CreateInt(func(observer IntObserver) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		time.Sleep(10 * time.Millisecond)
		observer.Next(3)
		observer.Complete()
	})

	sourceB := CreateInt(func(observer IntObserver) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})

	result, err := sourceA.Merge(sourceB).ToSlice()
	expect := []int{0, 1, 2, 3}

	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
