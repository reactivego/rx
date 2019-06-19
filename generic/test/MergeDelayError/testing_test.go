package MergeDelayError

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeDelayError(t *testing.T) {
	sourceA := CreateInt(func(observer IntObserver) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		observer.Error(RxError("error.sourceA"))
	})

	sourceB := CreateInt(func(observer IntObserver) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})

	result, err := sourceA.MergeDelayError(sourceB).ToSlice()
	expect := []int{0, 1, 2}

	assert.EqualError(t, err, "error.sourceA")
	assert.Equal(t, expect, result)
}
