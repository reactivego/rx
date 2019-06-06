package Throw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThrow(t *testing.T) {
	result, err := ThrowInt(RxError("throw")).ToSlice()
	expect := []int(nil)
	assert.EqualError(t, err, "throw")
	assert.Equal(t, expect, result)
}
