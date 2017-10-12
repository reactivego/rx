package ElementAt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElementAt(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4).ElementAt(2).ToSlice()
	expect := []int{3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
