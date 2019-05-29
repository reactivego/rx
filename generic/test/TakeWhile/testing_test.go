package TakeWhile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func smallerThan(value int) func(int) bool {
	return func(next int) bool { return next < value }
}

func TestTakeWhile(t *testing.T) {
	source := FromInts(1, 2, 3, 4, 5)

	result, err := source.TakeWhile(smallerThan(3)).ToSlice()
	expect := []int{1, 2}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	result, err = source.TakeWhile(smallerThan(4)).ToSlice()
	expect = []int{1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
