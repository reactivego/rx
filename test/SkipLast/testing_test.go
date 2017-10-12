package SkipLast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipLast(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4, 5).SkipLast(2).ToSlice()
	expect := []int{1, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
