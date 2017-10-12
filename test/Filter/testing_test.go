package Filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	even := func(i int) bool { return i%2 == 0 }
	result, err := FromInts(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).ToSlice()
	expect := []int{2, 4, 6, 8}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
