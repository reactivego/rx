package First

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4).First().ToSlice()
	expect := []int{1}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
