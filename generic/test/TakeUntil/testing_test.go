package TakeUntil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTakeUntil(t *testing.T) {
	interrupt := Never().Timeout(150 * time.Millisecond).Catch(Just("stop"))
	result, err := Interval(100 * time.Millisecond).TakeUntil(interrupt).ToSlice()

	expect := []int{0}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
