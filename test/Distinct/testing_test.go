package Distinct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinct(t *testing.T) {
	result, err := FromInts(1, 1, 2, 2, 3, 2, 4, 5).Distinct().ToSlice()
	expect := []int{1, 2, 3, 4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
