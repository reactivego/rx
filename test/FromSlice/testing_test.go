package FromSlice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	expect := []int{1, 2, 3, 4, 5}
	result, err := FromSliceInt(expect).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
