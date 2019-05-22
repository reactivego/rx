package Average

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverageInt(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4, 5).Average().ToSingle()
	expect := 3
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func TestAverageFloat32(t *testing.T) {
	result, err := FromFloat32s(1, 2, 3, 4).Average().ToSingle()
	expect := float32(2.5)
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
