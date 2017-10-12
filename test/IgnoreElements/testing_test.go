package IgnoreElements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreElements(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4, 5).IgnoreElements().ToSlice()
	expect := []int(nil)
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
