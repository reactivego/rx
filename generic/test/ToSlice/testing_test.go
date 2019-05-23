package ToSlice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSlice(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b, e := FromSliceInt(a).ToSlice()
	assert.NoError(t, e)
	assert.Equal(t, a, b)
}

func TestToSlice_twice(t *testing.T) {
	must := func(s []int, e error) []int {
		assert.NoError(t, e)
		return s
	}
	expected := []int{1, 2, 3, 4}
	actual := FromSliceInt(expected)
	assert.Equal(t, must(actual.ToSlice()), expected)
	assert.Equal(t, must(actual.ToSlice()), expected)
}
