package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ToSlice will run the observable to completion and then return the resulting
// slice and any error that was emitted by the observable. ToSlice will
// internally run the observable on an asynchronous scheduler.
func Example_toSlice() {
	if slice, err := Range(1, 9).ToSlice(); err == nil {
		for _, value := range slice {
			fmt.Print(value)
		}
	}

	// Output:
	// 123456789
}

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
