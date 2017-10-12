package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ToChan is used to convert the observable to a channel. ToChan changes the
// scheduler it uses for subscribing to an asynchronous one, so the subscription
// runs concurrently with the channel reading for loop.
func Example_toChan() {
	for value := range Range(1, 9).ToChan() {
		fmt.Print(value)
	}

	// Output:
	// 123456789
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
