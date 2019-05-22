package Skip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4, 5).Skip(2).ToSlice()
	expect := []int{3, 4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
