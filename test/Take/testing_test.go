package Take

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTake(t *testing.T) {
	source := FromInts(1, 2, 3, 4, 5)

	result, err := source.Take(2).ToSlice()
	expect := []int{1, 2}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	result, err = source.Take(3).ToSlice()
	expect = []int{1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
