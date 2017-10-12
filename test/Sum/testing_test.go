package Sum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumInt(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4, 5).Sum().ToSingle()
	expect := 15
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func TestSumFloat32(t *testing.T) {
	result, err := FromFloat32(1, 2, 3, 4.5).Sum().ToSingle()
	expect := float32(10.5)
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
