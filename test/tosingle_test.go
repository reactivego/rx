package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ToSingle is used to make sure only a single value was produced by the
// observable. ToSingle will internally run the observable on an asynchronous
// scheduler. However it will only return when the observable was complete or
// an error was emitted.
func Example_toSingle() {
	if value, err := FromInt(19).ToSingle(); err == nil {
		fmt.Println(value)
	}

	if _, err := FromInts(19, 20).ToSingle(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 19
	// expected one value, got multiple
}

func TestToSingle_errMultipleValue(t *testing.T) {
	_, err := FromInts(1, 2).ToSingle()
	assert.Error(t, err)
}

func TestToSingle_errNoValue(t *testing.T) {
	_, err := EmptyInt().ToSingle()
	assert.Error(t, err)
}

func TestToSingle_correct(t *testing.T) {
	value, err := FromInts(3).ToSingle()
	assert.NoError(t, err)
	assert.Equal(t, 3, value)
}
