package TakeLast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeLast(t *testing.T) {
	source := FromInts(1, 2, 3, 4, 5)

	result, err := source.TakeLast(2).ToSlice()
	expect := []int{4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	result, err = source.TakeLast(3).ToSlice()
	expect = []int{3, 4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
