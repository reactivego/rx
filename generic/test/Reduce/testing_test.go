package Reduce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	add := func(acc float32, value int) float32 {
		return acc + float32(value)
	}
	result, err := FromInts(1, 2, 3, 4, 5).ReduceFloat32(add, 0.0).ToSingle()
	expect := float32(15.0)
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
