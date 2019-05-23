package ToChan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToChanAny(t *testing.T) {
	expect := []interface{}{1, 2, 3, 4, 5, 4, 3, 2, 1}
	source := FromSlice(expect).ToChan()
	result := []interface{}{}
	var err error
	for next := range source {
		switch v := next.(type) {
		case int:
			result = append(result, v)
		case error:
			err = v
		}
	}
	assert.Equal(t, expect, result)
	assert.Nil(t, err)
}

func TestToChan(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	a := FromSliceInt(expected).ToChan()
	b := []int{}
	for i := range a {
		b = append(b, i)
	}
	assert.Equal(t, expected, b)
}
