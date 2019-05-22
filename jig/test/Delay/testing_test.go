package Delay

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelay(t *testing.T) {
	start := time.Now()
	actual, err := FromInts(1,2,3).Delay(time.Millisecond * 250).ToSlice()
	expect := []int{1,2,3}

	assert.NoError(t,err)
	assert.WithinDuration(t, start.Add(250*time.Millisecond), time.Now(), 10*time.Millisecond, "elapsed time must be between 245 and 255 ms")
	assert.Equal(t, expect, actual)
}
