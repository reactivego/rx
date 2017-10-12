package Range

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	source := Range(1, 5)
	result, err := source.ToSlice()
	expect := []int{1, 2, 3, 4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
	result, err = source.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
