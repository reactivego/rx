package From

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	result, err := FromInt(1, 2, 3, 4, 5).ToSlice()
	expect := []int{1, 2, 3, 4, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
